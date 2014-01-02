package main

import (
	"time"
)

type OrderType int

const (
	// marketOrder is defined as 0
	MarketOrder OrderType = iota

	// limitOrder is defined as 1
	LimitOrder
)

func (x OrderType) String() string {
	switch x {
	case MarketOrder:
		return "market"
	case LimitOrder:
		return "limit"
	}

	return ""
}

type OrderSide int

const (
	// iota is reset to 0, so bidSide is set to 0
	BidSide OrderSide = iota

	// askSide is set to 1
	AskSide
)

func (x OrderSide) String() string {
	switch x {
	case AskSide:
		return "ask"
	case BidSide:
		return "bid"
	}

	return ""
}

type OrderStatus int

const (
	PendingStatus OrderStatus = iota
	OpenStatus
	PartiallyFilledStatus
	CompletedStatus
	CanceledStatus
)

func (x OrderStatus) String() string {
	switch x {
	case PendingStatus:
		return "pending"
	case OpenStatus:
		return "open"
	case PartiallyFilledStatus:
		return "partially_filled"
	case CompletedStatus:
		return "completed"
	case CanceledStatus:
		return "canceled"
	}

	return ""
}

type Order struct {
	Uuid        string
	MarketUuid  string
	Size        int64
	InitialSize int64

	// price is not applicable when order is a marketOrder
	// in these cases price is a nil pointer
	Price *int64

	Side      OrderSide
	Status    OrderStatus
	Type      OrderType
	CreatedAt time.Time
}
