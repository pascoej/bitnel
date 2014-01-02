package main

import (
	"testing"
)

var expectedOrderType = map[OrderType]string{
	MarketOrder: "market",
	LimitOrder:  "limit",
}

var expectedOrderStatus = map[OrderStatus]string{
	PendingStatus:         "pending",
	OpenStatus:            "open",
	PartiallyFilledStatus: "partially_filled",
	CompletedStatus:       "completed",
	CanceledStatus:        "canceled",
}

var expectedOrderSide = map[OrderSide]string{
	BidSide: "bid",
	AskSide: "ask",
}

func TestOrderTypeString(t *testing.T) {
	for ty, s := range expectedOrderType {
		ts := ty.String()
		if ts != s {
			t.Errorf("%s expected, got %s", s, ts)
		}
	}
}

func TestOrderStatusString(t *testing.T) {
	for ty, s := range expectedOrderStatus {
		ts := ty.String()
		if ts != s {
			t.Errorf("%s expected, got %s", s, ts)
		}
	}
}

func TestOrderSideString(t *testing.T) {
	for ty, s := range expectedOrderSide {
		ts := ty.String()
		if ts != s {
			t.Errorf("%s expected, got %s", s, ts)
		}
	}
}
