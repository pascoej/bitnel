package model

import (
	"errors"
	"time"

	"code.google.com/p/go.crypto/bcrypt"
	vv "github.com/bitnel/bitnel/api/validator"
)

// Users have orders and identifying information. To place an order a user must
// be authenticated and authorized (will come in later).
type User struct {
	Uuid string `schema:"-" json:"uuid"`

	// email and password are the 2 fields that can be nil
	Email        string    `schema:"email" json:"email"`
	Password     string    `schema:"password" json:"-"`
	PasswordHash []byte    `schema:"-" json:"-"`
	CreatedAt    time.Time `schema:"-" json:"created_at"`
}

func (u *User) HashPassword(cost int) error {
	var err error

	u.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(u.Password), cost)
	if err != nil {
		return errors.New("model: error hashing user password")
	}

	return nil
}

func (u *User) ComparePassword(pass string) bool {
	// CompareHashAndPassword returns nil on success
	return nil == bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(pass))
}

func (u *User) Rules() map[string][]vv.Rule {
	return map[string][]vv.Rule{
		"Email":    []vv.Rule{&vv.NonZero{}, &vv.Email{}},
		"Password": []vv.Rule{&vv.NonZero{}, &vv.Length{3, 25}},
	}
}
