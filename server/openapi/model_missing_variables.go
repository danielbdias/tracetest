/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.2.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type MissingVariables struct {
	Key string `json:"key,omitempty"`

	DefaultValue string `json:"defaultValue,omitempty"`
}

// AssertMissingVariablesRequired checks if the required fields are not zero-ed
func AssertMissingVariablesRequired(obj MissingVariables) error {
	return nil
}

// AssertRecurseMissingVariablesRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of MissingVariables (e.g. [][]MissingVariables), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseMissingVariablesRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aMissingVariables, ok := obj.(MissingVariables)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertMissingVariablesRequired(aMissingVariables)
	})
}
