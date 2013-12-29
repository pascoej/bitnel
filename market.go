package main

type market struct {
	baseCurrency  *currency
	quoteCurrency *currency
}

func (m *market) currencyPair() string {
	return m.baseCurrency.code + m.quoteCurrency.code
}
