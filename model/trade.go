package model

import (
	"github.com/bitnel/bitnel-api/money"
	"time"
)

type Trade struct {
	Uuid               string
	OrderUuid          string
	AccountUuid        string
	TransactionUuid    string
	FeeTransactionUuid string
	CreatedAt          time.Time
}
