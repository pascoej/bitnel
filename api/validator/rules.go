package validator

import (
	"fmt"
	"strings"
)

var rules = map[string]Rule{
	"length": &Length{},
}

type Rule interface {
	errorMsg() string
	validate(interface{}) bool
}

type NonZero struct{}

func (v *NonZero) validate(t interface{}) bool {
	if t != nil {
		return true
	}

	return false
}

func (v *NonZero) errorMsg() string {
	return fmt.Sprint("is not a nonzero value")
}

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

func (v *Length) errorMsg() string {
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

func (v *MinLength) errorMsg() string {
	return fmt.Sprintf("does not meet length requirements (min: %d)", v.Min)
}

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

func (v *Set) errorMsg() string {
	return fmt.Sprintf("is not in the set %s", v.Set)
}

type Email struct{}

func (v *Email) validate(i interface{}) bool {
	if str, ok := i.(string); ok {
		if length := len(str); strings.Contains(str, "@") && (length >= 3 && length <= 254) {
			return true
		}
	}

	return false
}

func (v *Email) errorMsg() string {
	return fmt.Sprintf("is not a valid email")
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
