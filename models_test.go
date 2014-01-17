package main

import (
	"testing"
)

// Currency tests
var currencyTests = []struct {
	curr     Currency
	exString string
}{
	{Btc, "btc"},
	{Ltc, "ltc"},
}

func TestCurrency(t *testing.T) {
	for _, tt := range currencyTests {
		// test tt.curr.String() is equal to expected result
		if tt.curr.String() != tt.exString {
			t.Errorf("%s got, expected %s", tt.curr.String(), tt.exString)
		}
	}
}

// Order tests
var orderTypeTests = []struct {
	otype    OrderType
	exString string
}{
	{MarketOrder, "market"},
	{LimitOrder, "limit"},
}

func TestOrderType(t *testing.T) {
	for _, tt := range orderTypeTests {
		if tt.otype.String() != tt.exString {
			t.Error("%s expected, got %s", tt.exString, tt.otype.String())
		}
	}
}

var orderStatusTests = []struct {
	ostatus  OrderStatus
	exString string
}{
	{PendingStatus, "pending"},
	{OpenStatus, "open"},
	{PartiallyFilledStatus, "partially_filled"},
	{CompletedStatus, "completed"},
	{CanceledStatus, "canceled"},
}

func TestOrderStatus(t *testing.T) {
	for _, tt := range orderStatusTests {
		if tt.ostatus.String() != tt.exString {
			t.Error("%s expected, got %s", tt.exString, tt.ostatus.String())
		}
	}
}

var orderSideTests = []struct {
	oside    OrderSide
	exString string
}{
	{BidSide, "bid"},
	{AskSide, "ask"},
}

func TestOrderSide(t *testing.T) {
	for _, tt := range orderSideTests {
		if tt.oside.String() != tt.exString {
			t.Error("%s expected, got %s", tt.exString, tt.oside.String())
		}
	}
}

// User tests
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
