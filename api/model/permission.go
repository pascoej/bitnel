package model

import (
	"errors"
)

type Permission int

const (
	ViewBalances Permission = iota
	Addresses
	SubmitOrder
	ViewOrders
	Transactions
	OauthApps
	ViewUserInfo
)

func (x Permission) String() string {
	switch x {
		case ViewBalances: 
			return "view_balances"
		case Addresses:
			return "addresses"
		case SubmitOrder:
			return "submit_order"
		case ViewOrders:
			return "view_orders"
		case Transactions:
			return "transactions"
		case OauthApps:
			return "ouath_apps"
		case ViewUserInfo:
			return "view_user_info"
	}
	return ""
}

func ParsePermission(s string) (Permission, error) {
	switch s {
		case ViewBalances.String():
			return ViewBalances, nil
		case Addresses.String():
			return Addresses, nil
		case SubmitOrder.String():
			return SubmitOrder, nil
		case ViewOrders.String():
			return ViewOrders, nil
		case Transactions.String():
			return Transactions, nil
		case OauthApps.String():
			return OauthApps, nil
		case ViewUserInfo.String():
			return ViewUserInfo, nil
		default:
			return 0, errors.New("model: invalid Permission " + s)
	}
}