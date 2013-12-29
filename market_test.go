package main

import (
	"testing"
)

var testMarket = &market{
	baseCurrency:  ltc,
	quoteCurrency: btc,
}

func TestMarketCurrencyPair(t *testing.T) {
	if testMarket.currencyPair() != "ltcbtc" {
		t.Error("base currency is not ltcbtc")
	}
}
