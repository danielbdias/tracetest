package executor

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type traceDBFactoryFn func(ds model.DataStore) (tracedb.TraceDB, error)

type DefaultPollerExecutor struct {
	updater           RunUpdater
	newTraceDBFn      traceDBFactoryFn
	dsRepo            model.DataStoreRepository
	maxTracePollRetry int
}

type InstrumentedPollerExecutor struct {
	tracer         trace.Tracer
	pollerExecutor PollerExecutor
}

func (pe InstrumentedPollerExecutor) ExecuteRequest(request *PollingRequest) (bool, model.Run, error) {
	_, span := pe.tracer.Start(request.ctx, "Fetch trace")
	defer span.End()

	finished, run, err := pe.pollerExecutor.ExecuteRequest(request)

	spanCount := 0
	if run.Trace != nil {
		spanCount = len(run.Trace.Flat)
	}

	attrs := []attribute.KeyValue{
		attribute.String("tracetest.run.trace_poller.trace_id", request.run.TraceID.String()),
		attribute.String("tracetest.run.trace_poller.span_id", request.run.SpanID.String()),
		attribute.Bool("tracetest.run.trace_poller.succesful", finished),
		attribute.String("tracetest.run.trace_poller.test_id", string(request.test.ID)),
		attribute.Int("tracetest.run.trace_poller.amount_retrieved_spans", spanCount),
	}

	if err != nil {
		attrs = append(attrs, attribute.String("tracetest.run.trace_poller.error", err.Error()))
		span.RecordError(err)
	}

	span.SetAttributes(attrs...)
	return finished, run, err
}

func NewPollerExecutor(
	retryDelay time.Duration,
	maxWaitTimeForTrace time.Duration,
	tracer trace.Tracer,
	updater RunUpdater,
	newTraceDBFn traceDBFactoryFn,
	dsRepo model.DataStoreRepository,
) PollerExecutor {

	maxTracePollRetry := int(math.Ceil(float64(maxWaitTimeForTrace) / float64(retryDelay)))
	pollerExecutor := &DefaultPollerExecutor{
		updater:           updater,
		newTraceDBFn:      newTraceDBFn,
		dsRepo:            dsRepo,
		maxTracePollRetry: maxTracePollRetry,
	}

	return &InstrumentedPollerExecutor{
		tracer:         tracer,
		pollerExecutor: pollerExecutor,
	}
}

func (pe DefaultPollerExecutor) traceDB(ctx context.Context) (tracedb.TraceDB, error) {
	ds, err := pe.dsRepo.DefaultDataStore(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get default datastore: %w", err)
	}

	tdb, err := pe.newTraceDBFn(ds)
	if err != nil {
		return nil, fmt.Errorf(`cannot get tracedb from DataStore config with ID "%s": %w`, ds.ID, err)
	}

	return tdb, nil
}

func (pe DefaultPollerExecutor) ExecuteRequest(request *PollingRequest) (bool, model.Run, error) {
	log.Printf("[PollerExecutor] Test %s Run %d: ExecuteRequest\n", request.test.ID, request.run.ID)
	run := request.run

	traceDB, err := pe.traceDB(request.ctx)
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: GetDataStore error: %s\n", request.test.ID, request.run.ID, err.Error())
		return false, model.Run{}, err
	}

	traceID := run.TraceID.String()
	trace, err := traceDB.GetTraceByID(request.ctx, traceID)
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: GetTraceByID (traceID %s) error: %s\n", request.test.ID, request.run.ID, traceID, err.Error())
		return false, model.Run{}, err
	}

	trace.ID = run.TraceID

	if !pe.donePollingTraces(request, traceDB, trace) {
		log.Printf("[PollerExecutor] Test %s Run %d: Not done polling\n", request.test.ID, request.run.ID)
		run.Trace = &trace
		request.run = run
		return false, run, nil
	}

	log.Printf("[PollerExecutor] Test %s Run %d: Start Sorting\n", request.test.ID, request.run.ID)
	trace = trace.Sort()
	log.Printf("[PollerExecutor] Test %s Run %d: Sorting complete\n", request.test.ID, request.run.ID)
	run.Trace = &trace
	request.run = run

	if !trace.HasRootSpan() {
		newRoot := model.NewTracetestRootSpan(run)
		run.Trace = run.Trace.InsertRootSpan(newRoot)
	} else {
		run.Trace.RootSpan = model.AugmentRootSpan(run.Trace.RootSpan, run.TriggerResult)
	}
	run = run.SuccessfullyPolledTraces(run.Trace)

	fmt.Printf("completed polling result %d after %d times, number of spans: %d \n", run.ID, request.count, len(run.Trace.Flat))

	log.Printf("[PollerExecutor] Test %s Run %d: Start updating\n", request.test.ID, request.run.ID)
	err = pe.updater.Update(request.ctx, run)
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: Update error: %s\n", request.test.ID, request.run.ID, err.Error())
		return false, model.Run{}, err
	}

	return true, run, nil
}

func (pe DefaultPollerExecutor) donePollingTraces(job *PollingRequest, traceDB tracedb.TraceDB, trace model.Trace) bool {
	if !traceDB.ShouldRetry() {
		log.Printf("[PollerExecutor] Test %s Run %d: Done polling. TraceDB is not retryable\n", job.test.ID, job.run.ID)
		return true
	}
	// we're done if we have the same amount of spans after polling or `maxTracePollRetry` times
	if job.count == pe.maxTracePollRetry {
		log.Printf("[PollerExecutor] Test %s Run %d: Done polling. Hit MaxRetry of %d\n", job.test.ID, job.run.ID, pe.maxTracePollRetry)
		return true
	}

	if job.run.Trace == nil {
		return false
	}

	if len(trace.Flat) > traceDB.MinSpanCount() && len(trace.Flat) == len(job.run.Trace.Flat) {
		log.Printf("[PollerExecutor] Test %s Run %d: Done polling. Condition met\n", job.test.ID, job.run.ID)
		return true
	}

	return false
}
