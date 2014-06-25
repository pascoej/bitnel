package validator

import "reflect"

type Validator struct {
	data   Ruler
	errors []*FieldError
}

// models satifies
type Ruler interface {
	Rules() map[string][]Rule
}

type FieldError struct {
	Name string
	Val  interface{}
	Msg  string
}

func New(data Ruler) *Validator {
	return &Validator{data, []*FieldError{}}
}

// when fields is nil, the validator validates everything; on empty slice, it validates nothing
func (v *Validator) Validate(fields []string) (bool, []*FieldError) {
	rules := v.data.Rules()
	sv := reflect.ValueOf(v.data).Elem()

	if fields == nil {
		for fieldName, ruleSet := range rules {
			for _, rul := range ruleSet {
				cont := true // we validate until first rule fails
				if cont {
					dd := sv.FieldByName(fieldName).Interface()
					if val := rul.validate(dd); !val {
						ferr := &FieldError{
							Name: fieldName,
							Val:  dd,
							Msg:  rul.errorMsg(),
						}
						v.errors = append(v.errors, ferr)
						cont = false
					}
				}
			}
		}
	} else if len(fields) > 0 {
		for _, f := range fields {
			ruleSet, ok := rules[f]
			if !ok {
				panic("validator: field needs to be validated, but has no rules; check spelling")
			}

			for _, rul := range ruleSet {
				cont := true // we validate until first rule fails
				if cont {
					dd := sv.FieldByName(f).Interface()
					if val := rul.validate(dd); !val {
						ferr := &FieldError{
							Name: f,
							Val:  val,
							Msg:  rul.errorMsg(),
						}
						v.errors = append(v.errors, ferr)
						cont = false
					}
				}
			}
		}
	}

	return !(len(v.errors) > 0), v.errors
}
