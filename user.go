package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"time"
)

type User struct {
	Uuid         string
	Email        string
	Password     string
	PasswordHash []byte
	CreatedAt    time.Time
}

func (u *User) HashPassword(cost int) error {
	var err error

	u.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(u.Password), cost)
	if err != nil {
		return errors.New("error hashing user password")
	}

	return nil
}

func (u *User) ComparePassword(pass string) bool {
	// CompareHashAndPassword returns nil on success
	return nil == bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(pass))
}
