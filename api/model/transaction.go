package model

import (
	"github.com/bitnel/api/money"
	"time"
)

type TransactionType int

const (
	AdjustmentTransaction TransactionType = iota
	TradeTransaction
	WithdrawTransaction
	DepositTransaction
	OrderTransaction
)

func (t TransactionType) String() string {
	switch t {
	case AdjustmentTransaction:
		return "adjustment"
	case TradeTransaction:
		return "trade"
	case WithdrawTransaction:
		return "deposit"
	case DepositTransaction:
		return "withdraw"
	case OrderTransaction:
		return "order"
	}

	return ""
}

func ParseTransactionType(tt string) TransactionType {
	switch tt {
	case AdjustmentTransaction.String():
		return AdjustmentTransaction
	case TradeTransaction.String():
		return TradeTransaction
	case WithdrawTransaction.String():
		return WithdrawTransaction
	case DepositTransaction.String():
		return DepositTransaction
	case OrderTransaction.String():
		return OrderTransaction
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
	Trade       *string         `json:-`
	OrderUuid   *string         `json:-`
}

func (trans *Transaction) GetAmountAfterFee() money.Unit {
	return trans.Amount - trans.FeeAmount
}
