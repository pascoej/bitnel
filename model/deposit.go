package model

import (
	"github.com/bitnel/bitnel-api/money"
	"time"
)

type Deposit struct {
	Uuid            string     `json:"uuid"`
	TransactionUuid string     `json:"transaction_uuid"`
	AccountUuid     string     `json:"account_uuid"`
	Amount          money.Unit `json:"amount"`
	CreatedAt       time.Time  `json:"created_at"`
}
