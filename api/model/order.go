package model

import (
	//"database/sql"
	"errors"
	//"fmt"
	"github.com/bitnel/api/money"
	"time"
)

// // Order type defines the order's type. Currently the exchange only supports
// // market and limit orders.
// type OrderType int

// const (
// 	// marketOrder is defined as 0
// 	MarketOrder OrderType = iota

// 	// limitOrder is defined as 1
// 	LimitOrder
// )

// func (x OrderType) String() string {
// 	switch x {
// 	case MarketOrder:
// 		return "market"
// 	case LimitOrder:
// 		return "limit"
// 	}

// 	return ""
// }

// func ParseOrderType(s string) (OrderType, error) {
// 	switch s {
// 	case MarketOrder.String():
// 		return MarketOrder, nil
// 	case LimitOrder.String():
// 		return LimitOrder, nil
// 	default:
// 		// 0 because underlying type of OrderType is int, and int "zero value" is 0
// 		return 0, errors.New("model: invalid OrderType " + s)
// 	}
// }

// func (x OrderType) MarshalJSON() ([]byte, error) {
// 	return []byte("\"" + x.String() + "\""), nil
// }

// func (x *OrderType) UnmarshalJSON(b []byte) error {
// 	var err error
// 	if *x, err = ParseOrderType(string(b[1 : len(b)-1])); err != nil {
// 		return err
// 	}

// 	return nil
// }

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

func ParseOrderSide(s string) (OrderSide, error) {
	//fmt.Println(s)
	switch s {
	case AskSide.String():
		return AskSide, nil
	case BidSide.String():
		return BidSide, nil
	default:
		return 0, errors.New("model: invalid OrderSide " + s)
	}
}

func (x OrderSide) CounterSide() OrderSide {
	if x == BidSide {
		return AskSide
	}
	return BidSide
}

func (x OrderSide) MarshalJSON() ([]byte, error) {
	return []byte("\"" + x.String() + "\""), nil
}

func (x *OrderSide) UnmarshalJSON(b []byte) error {
	var err error
	if *x, err = ParseOrderSide(string(b[1 : len(b)-1])); err != nil {
		return err
	}

	return nil
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

func ParseOrderStatus(s string) (OrderStatus, error) {
	switch s {
	case PendingStatus.String():
		return PendingStatus, nil
	case OpenStatus.String():
		return OpenStatus, nil
	case PartiallyFilledStatus.String():
		return PartiallyFilledStatus, nil
	case CompletedStatus.String():
		return CompletedStatus, nil
	case CanceledStatus.String():
		return CanceledStatus, nil
	default:
		return 0, errors.New("model: invalid OrderStatus " + s)
	}
}

func (x OrderStatus) MarshalJSON() ([]byte, error) {
	return []byte("\"" + x.String() + "\""), nil
}

func (x *OrderStatus) UnmarshalJSON(b []byte) error {
	var err error
	if *x, err = ParseOrderStatus(string(b[1 : len(b)-1])); err != nil {
		return err
	}

	return nil
}

// A order defines the information belonging to an order, which is also part of
// a market.
type Order struct {
	Uuid        string      `json:"uuid"`
	MarketUuid  *string     `json:"market_uuid"`
	Size        *money.Unit `json:"size"`
	InitialSize money.Unit  `json:"initial_size"`
	Price       *money.Unit `json:"price"`
	Side        *OrderSide  `json:"side"`
	Status      OrderStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	AccountUuid *string     `json:"-"`

	// Type        *OrderType  `json:"type"`
}
