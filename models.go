package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"time"
)

// A market defines what currency is being exchanged, and what currency is used
// for the quoting price.
type Market struct {
	Uuid          string
	BaseCurrency  Currency
	QuoteCurrency Currency

	// need to specify for url routing
	CurrencyPair string
}

// A currency defines a currency being used in the system.
type Currency int

const (
	Btc Currency = iota
	Ltc
)

func (c Currency) String() string {
	switch c {
	case Btc:
		return "btc"
	case Ltc:
		return "ltc"
	}

	return ""
}

// Order type defines the order's type. Currently the exchange only supports
// market and limit orders.
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

// Order side defines the side which an order lies. For example a buy order
// would mean the order is on the bid side.
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

// Order status defines the current status of an order.
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

// A order defines the information belonging to an order, which is also part of
// a market.
type Order struct {
	Uuid        string
	MarketUuid  string
	Size        *int64
	InitialSize *int64

	// price is not applicable when order is a marketOrder
	// in these cases price is a nil pointer
	Price *int64

	Side      OrderSide
	Status    OrderStatus
	Type      OrderType
	CreatedAt time.Time
}

// Users have orders and identifying information. To place an order a user must
// be authenticated and authorized (will come in later).
type User struct {
	Uuid         string
	Email        string
	Password     string
	PasswordHash []byte
	CreatedAt    time.Time
}

func (u *User) HashPassword(cost int) error {
	var err error

	u.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(u.Password), cost)
	if err != nil {
		return errors.New("error hashing user password")
	}

	return nil
}

func (u *User) ComparePassword(pass string) bool {
	// CompareHashAndPassword returns nil on success
	return nil == bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(pass))
}
