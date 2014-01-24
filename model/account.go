package model

import (
	"github.com/bitnel/bitnel-api/money"
)

type Account struct {
	Uuid             string         `json:"uuid"`
	UserUuid         string         `json:"uuid"`
	Currency         money.Currency `json:"currency"`
	AvailableBalance money.Unit     `json:"available_balance"`
	ReservedBalance  money.Unit     `json:"reserved_balance"`
}

func (a *Account) NetBalance() money.Unit {
	return a.AvailableBalance + a.ReservedBalance
}
