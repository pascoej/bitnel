package main

type market struct {
	uuid          string
	baseCurrency  currency
	quoteCurrency currency

	// need to specify for url routing
	currencyPair string
}
