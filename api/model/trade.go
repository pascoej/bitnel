package model

import (
	"github.com/bitnel/api/money"
	"time"
)

type Trade struct {
	Uuid      string     `json:"uuid"`
	Amount    money.Unit `json:"amount"`
	Price     money.Unit `json:"price"`
	CreatedAt time.Time  `json:"created_at"`
}
