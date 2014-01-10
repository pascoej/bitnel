package validation

import (
	"fmt"
	"strings"
)

// All validation rules should implement these methods
type Rule interface {
	errMsg() string

	// Validator should return true with whitelisting. False on default.
	validate(interface{}) bool
}

type Required struct{}

// Golang is a pain with zero values. When unmarshaling JSON to a struct with a bool field
// false may mean no value was provided or the field is _actually_ false.
func (v *Required) validate(t interface{}) bool {
	if t != nil {
		return true
	}

	return false
}

func (v *Required) errMsg() string {
	return "is not provided"
}

// Length validator checks the length of a string
type Length struct {
	Min int
	Max int
}

func (v *Length) validate(i interface{}) bool {
	if str, ok := i.(string); ok {
		if length := len(str); length >= v.Min && length <= v.Max {
			return true
		}
	}

	return false
}

func (v *Length) errMsg() string {
	return fmt.Sprintf("does not meet length requirements (min: %d, max: %d)", v.Min, v.Max)
}

// Mininum length validator
type MinLength struct {
	Min int
}

func (v *MinLength) validate(i interface{}) bool {
	if str, ok := i.(string); ok {
		if length := len(str); length >= v.Min {
			return true
		}
	}

	return false
}

func (v *MinLength) errMsg() string {
	return fmt.Sprintf("does not meet length requirements (min: %d)", v.Min)
}

// Set validator. Ensures a value is in a set
type Set struct {
	Set []string
}

func (v *Set) validate(i interface{}) bool {
	if str, ok := i.(string); ok {
		for _, val := range v.Set {
			if str == val {
				return true
			}
		}
	}

	return false
}

func (v *Set) errMsg() string {
	return fmt.Sprintf("is not in the set %s", v.Set)
}

// Email validator which checks for @ and length, real validation should be done
// by sending an confirmation email, duh.
type Email struct{}

func (v *Email) validate(i interface{}) bool {
	if str, ok := i.(string); ok {
		if length := len(str); strings.Contains(str, "@") && (length >= 3 && length <= 254) {
			return true
		}
	}

	return false
}

func (v *Email) errMsg() string {
	return "does not pass email validation"
}

// Min validator
type Min struct {
	Min int
}

func (v *Min) validate(i interface{}) bool {
	if num, ok := i.(int); ok {
		if num >= v.Min {
			return true
		}
	}

	return false
}

func (v *Min) errMsg() string {
	return fmt.Sprintf("is not greater than %d", v.Min)
}
