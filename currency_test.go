package main

import (
	"testing"
)

var expectedCurrencyCode = map[currency]string{
	btc: "btc",
	ltc: "ltc",
}

func TestCurrencyString(t *testing.T) {
	for ty, s := range expectedCurrencyCode {
		ts := ty.String()
		if ts != s {
			t.Errorf("%s expected, got %s", s, ts)
		}
	}
}
