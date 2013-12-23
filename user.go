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

func (u *user) hashPassword() error {
	var err error

	// Need to change cost later
	u.passwordHash, err = bcrypt.GenerateFromPassword([]byte(u.password), 8)
	if err != nil {
		return errors.New("error hashing user password")
	}

	return nil
}

func (u *user) comparePassword(pass string) bool {
	// CompareHashAndPassword returns nil on success
	return nil == bcrypt.CompareHashAndPassword(u.passwordHash, []byte(pass))
}
