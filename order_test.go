package main

import (
	"testing"
)

var expectedOrderType = map[orderType]string{
	marketOrder: "market",
	limitOrder:  "limit",
}

var expectedOrderStatus = map[orderStatus]string{
	pendingStatus:         "pending",
	openStatus:            "open",
	partiallyFilledStatus: "partially_filled",
	completedStatus:       "completed",
	canceledStatus:        "canceled",
}

var expectedOrderSide = map[orderSide]string{
	bidSide: "bid",
	askSide: "ask",
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
