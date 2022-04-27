/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"context"
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	openapi "github.com/kubeshop/tracetest/server/go"
	"github.com/kubeshop/tracetest/server/go/analytics"
	"github.com/kubeshop/tracetest/server/go/executor"
	"github.com/kubeshop/tracetest/server/go/testdb"
	"github.com/kubeshop/tracetest/server/go/tracedb"
	"github.com/kubeshop/tracetest/server/go/tracedb/jaegerdb"
	"github.com/kubeshop/tracetest/server/go/tracedb/tempodb"
	"github.com/kubeshop/tracetest/server/go/websocket"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var cfg = flag.String("config", "config.yaml", "path to the config file")

func main() {
	flag.Parse()
	c, err := LoadConfig(*cfg)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	tp := initOtelTracing(ctx)
	defer func() { _ = tp.Shutdown(ctx) }()

	testDB, err := testdb.New(c.PostgresConnString)
	if err != nil {
		log.Fatal(err)
	}

	var traceDB tracedb.TraceDB
	switch {
	case c.JaegerConnectionConfig != nil:
		log.Printf("connecting to Jaeger: %s\n", c.JaegerConnectionConfig.Endpoint)
		traceDB, err = jaegerdb.New(c.JaegerConnectionConfig)
		if err != nil {
			log.Fatal(err)
		}
	case c.TempoConnectionConfig != nil:
		log.Printf("connecting to tempo: %s\n", c.TempoConnectionConfig.Endpoint)
		traceDB, err = tempodb.New(c.TempoConnectionConfig)
		if err != nil {
			log.Fatal(err)
		}
	}

	ex, err := executor.New()
	if err != nil {
		log.Fatal(err)
	}

	maxWaitTimeForTrace, err := time.ParseDuration(c.MaxWaitTimeForTrace)
	if err != nil {
		// use a default value
		maxWaitTimeForTrace = 30 * time.Second
	}

	tracePoller := openapi.NewTracePoller(traceDB, testDB, maxWaitTimeForTrace)
	tracePoller.Start(5) // worker count. should be configurable
	defer tracePoller.Stop()

	runner := openapi.NewPersistentRunner(ex, testDB, tracePoller)
	runner.Start(5) // worker count. should be configurable
	defer runner.Stop()

	apiApiService := openapi.NewApiApiService(traceDB, testDB, runner)
	apiApiController := openapi.NewApiApiController(apiApiService)

	router := openapi.NewRouter(apiApiController)
	router.Use(otelmux.Middleware("tracetest"))

	dir := "./html"
	fileServer := http.FileServer(http.Dir(dir))
	fileMatcher := regexp.MustCompile(`\.[a-zA-Z]*$`)
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !fileMatcher.MatchString(r.URL.Path) {
			serveIndex(w, dir+"/index.html")
		} else {
			fileServer.ServeHTTP(w, r)
		}
	})

	err = analytics.CreateAndSendEvent("server_started", "beacon")
	if err != nil {
		log.Fatal(err)
	}

	go startWebsocketServer()
	log.Printf("HTTP Server started")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func startWebsocketServer() {
	wsRouter := websocket.NewRouter()
	wsRouter.Add("subscribe", websocket.HandleSubscribeCommand)
	log.Printf("WS Server started")

	wsRouter.ListenAndServe(":8081")
}

type gaParams struct {
	MeasurementId    string
	AnalyticsEnabled bool
}

func serveIndex(w http.ResponseWriter, path string) {
	templateData := gaParams{
		MeasurementId:    os.Getenv("GOOGLE_ANALYTICS_MEASUREMENT_ID"),
		AnalyticsEnabled: os.Getenv("ANALYTICS_ENABLED") == "true",
	}

	tpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err = tpl.Execute(w, templateData); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func initOtelTracing(ctx context.Context) *sdktrace.TracerProvider {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	var (
		exporter sdktrace.SpanExporter
		err      error
	)

	if endpoint == "" {
		endpoint = "opentelemetry-collector:4317"
		exporter, err = stdouttrace.New(stdouttrace.WithWriter(io.Discard))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(endpoint),
			otlptracegrpc.WithInsecure(),
		}
		exporter, err = otlptrace.New(ctx, otlptracegrpc.NewClient(opts...))
		if err != nil {
			log.Fatal(err)
		}
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{}))

	// Set standard attributes per semantic conventions
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("tracetest"),
	)

	// Create and set the TraceProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	return tp
}
