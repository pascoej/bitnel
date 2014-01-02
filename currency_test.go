package main

import (
	"testing"
)

var currencyTests = []struct {
	curr     Currency
	exString string
}{
	{Btc, "btc"},
	{Ltc, "ltc"},
}

func TestCurrency(t *testing.T) {
	for _, tt := range currencyTests {
		// test tt.curr.String() is equal to expected result
		if tt.curr.String() != tt.exString {
			t.Errorf("%s got, expected %s", tt.curr.String(), tt.exString)
		}
	}
}
