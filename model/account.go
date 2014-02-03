package model

type Account struct {
<<<<<<< HEAD
	Uuid     string      `json:"uuid"`
	UserUuid string      `json:"user_uuid"`
	Type     AccountType `json:"type"`
}
type AccountType int

const (
	accountTypeExchange AccountType = iota
)
=======
	Uuid     string
	UserUuid string
}
>>>>>>> 8822ead2d45d8caa6d290ab78fc0e24a8ef488d4
