package main

import (
	"time"
)

type orderType int
type orderSide int

const (
	// marketOrder is defined as 0
	marketOrder orderType = iota

	// limitOrder is defined as 1
	limitOrder

	// iota is reset to 0, so bidSide is set to 0
	bidSide orderSide = iota

	// askSide is set to 1
	askSide
)

type order struct {
	uuid         string
	size         int64
	initial_size int64

	// price is not applicable when order is a marketOrder
	// in these cases price is a nil pointer
	price *int64

	side orderSide

	// `type` is a reseved golang keyword
	oType orderType

	created_at time.Time
}
