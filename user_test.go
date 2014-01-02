package main

import (
	"testing"
)

var testBcryptCost = 10

// hashPassword() should hash password
func TestUserHashPassword(t *testing.T) {
	usr := &User{Password: "asdfasdf"}
	usr.HashPassword(testBcryptCost)

	if len(usr.PasswordHash) <= 0 {
		t.Error("PasswordHash should not be empty")
	}
}

// comparePassword() should compare correctly
func TestUserComparePassword(t *testing.T) {
	usr := &User{Password: "asdfasdf"}
	usr.HashPassword(testBcryptCost)

	if !usr.ComparePassword("asdfasdf") {
		t.Error("ComparePassword() should work")
	}
}
