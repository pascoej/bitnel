package model

import (
	"github.com/bitnel/bitnel-api/money"
	"time"
)

type TransactionType int

const (
	FeeTransaction TransactionType = iota
	OutboundTransaction
	InboundTransaction
)

func (t TransactionType) String() string {
	switch t {
	case FeeTransaction:
		return "fee"
	case OutboundTransaction:
		return "out"
	case InboundTransaction:
		return "in"
	}
}

func ParseTransactionType(tt string) TransactionType {
	switch tt {
	case FeeTransaction.String():
		return FeeTransaction
	case OutboundTransaction.String()
		return OutboundTransaction
	}
}

type Transaction struct {
	Uuid string `json:"uuid"`
	AccountUuid string  `json:"amount"`
	Type TransactionType `json:"amount"`
	Amount money.Unit `json:"amount"`
	CreatedAt time.Time `json:"amount"`
}
