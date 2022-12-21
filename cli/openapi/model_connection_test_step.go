/*
TraceTest

OpenAPI definition for TraceTest endpoint and resources

API version: 0.2.1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// ConnectionTestStep struct for ConnectionTestStep
type ConnectionTestStep struct {
	Passed  *bool   `json:"passed,omitempty"`
	Message *string `json:"message,omitempty"`
	Error   *string `json:"error,omitempty"`
}

// NewConnectionTestStep instantiates a new ConnectionTestStep object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConnectionTestStep() *ConnectionTestStep {
	this := ConnectionTestStep{}
	return &this
}

// NewConnectionTestStepWithDefaults instantiates a new ConnectionTestStep object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConnectionTestStepWithDefaults() *ConnectionTestStep {
	this := ConnectionTestStep{}
	return &this
}

// GetPassed returns the Passed field value if set, zero value otherwise.
func (o *ConnectionTestStep) GetPassed() bool {
	if o == nil || o.Passed == nil {
		var ret bool
		return ret
	}
	return *o.Passed
}

// GetPassedOk returns a tuple with the Passed field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectionTestStep) GetPassedOk() (*bool, bool) {
	if o == nil || o.Passed == nil {
		return nil, false
	}
	return o.Passed, true
}

// HasPassed returns a boolean if a field has been set.
func (o *ConnectionTestStep) HasPassed() bool {
	if o != nil && o.Passed != nil {
		return true
	}

	return false
}

// SetPassed gets a reference to the given bool and assigns it to the Passed field.
func (o *ConnectionTestStep) SetPassed(v bool) {
	o.Passed = &v
}

// GetMessage returns the Message field value if set, zero value otherwise.
func (o *ConnectionTestStep) GetMessage() string {
	if o == nil || o.Message == nil {
		var ret string
		return ret
	}
	return *o.Message
}

// GetMessageOk returns a tuple with the Message field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectionTestStep) GetMessageOk() (*string, bool) {
	if o == nil || o.Message == nil {
		return nil, false
	}
	return o.Message, true
}

// HasMessage returns a boolean if a field has been set.
func (o *ConnectionTestStep) HasMessage() bool {
	if o != nil && o.Message != nil {
		return true
	}

	return false
}

// SetMessage gets a reference to the given string and assigns it to the Message field.
func (o *ConnectionTestStep) SetMessage(v string) {
	o.Message = &v
}

// GetError returns the Error field value if set, zero value otherwise.
func (o *ConnectionTestStep) GetError() string {
	if o == nil || o.Error == nil {
		var ret string
		return ret
	}
	return *o.Error
}

// GetErrorOk returns a tuple with the Error field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectionTestStep) GetErrorOk() (*string, bool) {
	if o == nil || o.Error == nil {
		return nil, false
	}
	return o.Error, true
}

// HasError returns a boolean if a field has been set.
func (o *ConnectionTestStep) HasError() bool {
	if o != nil && o.Error != nil {
		return true
	}

	return false
}

// SetError gets a reference to the given string and assigns it to the Error field.
func (o *ConnectionTestStep) SetError(v string) {
	o.Error = &v
}

func (o ConnectionTestStep) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Passed != nil {
		toSerialize["passed"] = o.Passed
	}
	if o.Message != nil {
		toSerialize["message"] = o.Message
	}
	if o.Error != nil {
		toSerialize["error"] = o.Error
	}
	return json.Marshal(toSerialize)
}

type NullableConnectionTestStep struct {
	value *ConnectionTestStep
	isSet bool
}

func (v NullableConnectionTestStep) Get() *ConnectionTestStep {
	return v.value
}

func (v *NullableConnectionTestStep) Set(val *ConnectionTestStep) {
	v.value = val
	v.isSet = true
}

func (v NullableConnectionTestStep) IsSet() bool {
	return v.isSet
}

func (v *NullableConnectionTestStep) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConnectionTestStep(val *ConnectionTestStep) *NullableConnectionTestStep {
	return &NullableConnectionTestStep{value: val, isSet: true}
}

func (v NullableConnectionTestStep) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConnectionTestStep) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}