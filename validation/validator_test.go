package validation

import (
	"testing"
)

func TestValidator(t *testing.T) {
	var data = map[string]interface{}{
		"name":     "12381928312389128391823918274098127938172983712",
		"passwodf": "adjf9819237",
	}

	var vdr = New(data)

	vdr.Rule("password", &Required{})

	if vdr.Validate() == nil {
		t.Error("there should be errors returned")
	}

	vdr.Clear()

	if vdr.Validate() != nil {
		t.Error("clear does not work")
	}
}
