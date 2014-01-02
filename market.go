package main

type Market struct {
	Uuid          string
	BaseCurrency  Currency
	QuoteCurrency Currency

	// need to specify for url routing
	CurrencyPair string
}
