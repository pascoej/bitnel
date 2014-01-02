package main

import (
	"testing"
)

var expectedCurrencyCode = map[Currency]string{
	Btc: "btc",
	Ltc: "ltc",
}

func TestCurrencyString(t *testing.T) {
	for ty, s := range expectedCurrencyCode {
		ts := ty.String()
		if ts != s {
			t.Errorf("%s expected, got %s", s, ts)
		}
	}
}
