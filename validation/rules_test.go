package validation

import (
	"testing"
)

var ruleTests = []struct {
	rule Rule

	// There can be multiple pass/fail tests
	shouldPass []interface{}
	shouldFail []interface{}
}{
	{&Required{}, []interface{}{"asdf", 1232}, []interface{}{nil}},
	{&MinLength{5}, []interface{}{"12345", "123459"}, []interface{}{"123", "1234", nil, 123}},
	{&Length{2, 4}, []interface{}{"12", "1234", "123"}, []interface{}{"1", "12345", nil, 123}},
}

func TestRules(t *testing.T) {
	for _, tt := range ruleTests {
		// first do the passes
		for _, pp := range tt.shouldPass {
			if !tt.rule.validate(pp) {
				t.Errorf("value %s should not pass", pp)
			}
		}

		for _, pp := range tt.shouldFail {
			if tt.rule.validate(pp) {
				t.Errorf("value %s should not pass", pp)
			}
		}
	}
}
