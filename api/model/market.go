package model

import (
	"github.com/bitnel/api/money"
)

// A market defines what currency is being exchanged, and what currency is used
// for the quoting price.
type Market struct {
	Uuid          string
	BaseCurrency  money.Currency
	QuoteCurrency money.Currency

	// need to specify for url routing
	CurrencyPair string
}
