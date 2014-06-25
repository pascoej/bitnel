package model

import (
	"github.com/bitnel/api/money"
)

type Withdraw struct {
	Uuid            string `json:"uuid"`
	TransactionUuid string `json:"transaction_uuid"`
	AccountUuid     string `json:"account_uuid"`
	Size            money.Unit
}
