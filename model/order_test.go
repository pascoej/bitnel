package model

import (
	"encoding/json"
	"testing"
)

func TestParseOrderType(t *testing.T) {
	if l, err := ParseOrderType(MarketOrder.String()); err != nil || l != MarketOrder {
		t.Error("we have problem")
	}
}

func TestOrderTypeUnmarshalJSON(t *testing.T) {
	b := []byte(`{"type": "limit"}`)

	ff := &Order{}

	if err := json.Unmarshal(b, &ff); err != nil {
		t.Error("encountered an error while", err)
	}

	if ff.Type.String() != "limit" {
		t.Error("lol")
	}
}

func TestOrderSideUnmarshalJSON(t *testing.T) {
	b := []byte(`{"side": "bid"}`)

	ff := &Order{}

	if err := json.Unmarshal(b, &ff); err != nil {
		t.Error("encountered an error while", err)
	}

	if ff.Side.String() != "bid" {
		t.Error("lol")
	}
}

func TestOrderStatusUnmarshalJSON(t *testing.T) {
	b := []byte(`{"status": "canceled"}`)

	ff := &Order{}

	if err := json.Unmarshal(b, &ff); err != nil {
		t.Error("encountered an error while", err)
	}

	if ff.Status.String() != "canceled" {
		t.Error("lol")
	}
}
