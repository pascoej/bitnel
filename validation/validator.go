package validation

import (
	"fmt"
)

type ValidationError struct {
	errors *[]FieldError
}

// An error occured while validating a field
// Name should be the name of the field, and
// Reason should be the reason of the failure
type FieldError struct {
	Field string

	// Reasons should start with lowercase so that Field and Reason can be
	// concatenated with no grammatical errors
	Reason string
}

// Technically an error
func (f *FieldError) Error() string {
	return fmt.Sprintf("%s %s", f.Field, f.Reason)
}

// Stores errors occured in a validation session
type Validator struct {
	data   map[string]interface{}
	errors []*FieldError
}

// Starts a validation session
func New(d map[string]interface{}) *Validator {
	return &Validator{d, []*FieldError{}}
}

// Add a rule
func (v *Validator) Rule(key string, r Rule) {
	if !r.validate(v.data[key]) {
		v.errors = append(v.errors, &FieldError{key, r.errMsg()})
	}
}

// Returns nil if no errors, otherwise return the _first error_ for each field
func (v *Validator) Validate() map[string]*FieldError {
	if len(v.errors) > 0 {
		errors := map[string]*FieldError{}

		for _, e := range v.errors {
			if _, ok := errors[e.Field]; !ok {
				errors[e.Field] = e
			}
		}

		return errors
	}

	return nil
}

// Clears the validation session of any errors
func (v *Validator) Clear() {
	v.errors = []*FieldError{}
}
