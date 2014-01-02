package main

import (
	"testing"
)

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
