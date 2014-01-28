package model

import (
	"github.com/bitnel/bitnel-api/money"
	"time"
)

type TransactionType int

const (
	AdjustmentTransaction TransactionType = iota
	OutboundTransaction
	InboundTransaction
	WithdrawTransaction
	DepositTransaction
)

func (t TransactionType) String() string {
	switch t {
	case AdjustmentTransaction:
		return "adjustment"
	case OutboundTransaction:
		return "out"
	case InboundTransaction:
		return "in"
	case WithdrawTransaction:
		return "deposit"
	case DepositTransaction:
		return "withdraw"
	}

	return ""
}

func ParseTransactionType(tt string) TransactionType {
	switch tt {
	case AdjustmentTransaction.String():
		return AdjustmentTransaction
	case OutboundTransaction.String():
		return OutboundTransaction
	case InboundTransaction.String():
		return InboundTransaction
	case WithdrawTransaction.String():
		return WithdrawTransaction
	case DepositTransaction.String():
		return DepositTransaction
	}

	return 0
}

type Transaction struct {
	Uuid        string          `json:"uuid"`
	AccountUuid string          `json:"account_uuid"`
	Type        TransactionType `json:"type"`
	Amount      money.Unit      `json:"amount"`
	FeeAmount   money.Unit      `json:"fee_amount"`
	CreatedAt   time.Time       `json:"created_at"`
}
