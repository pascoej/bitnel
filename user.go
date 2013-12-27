package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"time"
)

type user struct {
	uuid         string
	email        string
	password     string
	passwordHash []byte
	created_at   time.Time
}

func (u *user) hashPassword(cost int) error {
	var err error

	u.passwordHash, err = bcrypt.GenerateFromPassword([]byte(u.password), cost)
	if err != nil {
		return errors.New("error hashing user password")
	}

	return nil
}

func (u *user) comparePassword(pass string) bool {
	// CompareHashAndPassword returns nil on success
	return nil == bcrypt.CompareHashAndPassword(u.passwordHash, []byte(pass))
}
