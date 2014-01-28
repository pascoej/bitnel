package model

import (
	"time"
)

type Trade struct {
	Uuid                  string    `json:"uuid"`
	OrderUuid             string    `json:"order_uuid"`
	AccountUuid           string    `json:"account_uuid"`
	BuyerTransactionUuid  string    `json:"buyer_transaction_uuid"`
	SellerTransactionUuid string    `json:"seller_transaction_uuid"`
	CreatedAt             time.Time `json:"created_at"`
}
