package model

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"time"
)

// Users have orders and identifying information. To place an order a user must
// be authenticated and authorized (will come in later).
type User struct {
	Uuid string `json:"uuid"`

	// email and password are the 2 fields that can be nil
	Email        *string   `json:"email"`
	Password     *string   `json:"password,omitempty"`
	PasswordHash []byte    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

func (u *User) HashPassword(cost int) error {
	var err error

	u.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(*u.Password), cost)
	if err != nil {
		return errors.New("model: error hashing user password")
	}

	return nil
}

func (u *User) ComparePassword(pass string) bool {
	// CompareHashAndPassword returns nil on success
	return nil == bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(pass))
}
