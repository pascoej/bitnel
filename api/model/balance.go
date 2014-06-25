package model

import (
	"github.com/bitnel/bitnel/api/money"
)

type Balance struct {
	Uuid             string         `json:"uuid"`
	UserUuid         string         `json:"uuid"`
	Currency         money.Currency `json:"currency"`
	AvailableBalance money.Unit     `json:"available_balance"`
	ReservedBalance  money.Unit     `json:"reserved_balance"`
}

func (a *Balance) NetBalance() money.Unit {
	return a.AvailableBalance + a.ReservedBalance
}
