package main

import (
	"testing"
)

var bcryptCost = 10

// hashPassword() should hash password
func TestUserHashPassword(t *testing.T) {
	usr := &user{password: "asdfasdf"}
	usr.hashPassword(bcryptCost)

	if len(usr.passwordHash) <= 0 {
		t.Error("passwordHash should not be empty")
	}
}

// comparePassword() should compare correctly
func TestUserComparePassword(t *testing.T) {
	usr := &user{password: "asdfasdf"}
	usr.hashPassword(bcryptCost)

	if !usr.comparePassword("asdfasdf") {
		t.Error("comparePassword() should work")
	}
}
