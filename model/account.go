package model

type Account struct {
	Uuid     string      `json:"uuid"`
	UserUuid string      `json:"user_uuid"`
	Type     AccountType `json:"type"`
}
type AccountType int

const (
	accountTypeExchange AccountType = iota
)
