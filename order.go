package main

import (
	"time"
)

type orderType int

const (
	// marketOrder is defined as 0
	marketOrder orderType = iota

	// limitOrder is defined as 1
	limitOrder
)

func (x orderType) String() string {
	switch x {
	case marketOrder:
		return "market"
	case limitOrder:
		return "limit"
	}

	return ""
}

type orderSide int

const (
	// iota is reset to 0, so bidSide is set to 0
	bidSide orderSide = iota

	// askSide is set to 1
	askSide
)

func (x orderSide) String() string {
	switch x {
	case askSide:
		return "ask"
	case bidSide:
		return "bid"
	}

	return ""
}

type orderStatus int

const (
	pendingStatus orderStatus = iota
	openStatus
	partiallyFilledStatus
	completedStatus
	canceledStatus
)

func (x orderStatus) String() string {
	switch x {
	case pendingStatus:
		return "pending"
	case openStatus:
		return "open"
	case partiallyFilledStatus:
		return "partially_filled"
	case completedStatus:
		return "completed"
	case canceledStatus:
		return "canceled"
	}

	return ""
}

type order struct {
	uuid         string
	size         int64
	initial_size int64

	// price is not applicable when order is a marketOrder
	// in these cases price is a nil pointer
	price *int64

	side   orderSide
	status orderStatus

	// `type` is a reseved golang keyword
	oType orderType

	created_at time.Time
}
