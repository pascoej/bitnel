package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitnel/bitnel-api/model"
	"github.com/bitnel/bitnel-api/money"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// This serves the root path of our API. Be friendly; say hello.
func indexHandler(w http.ResponseWriter, r *http.Request) *serverError {
	fmt.Fprintln(w, "Welcome to the Bitnel API!")

	return nil
}

// We define our own not found handler because we dislike the default Gorilla
// 404 message.
func notFoundHandler(w http.ResponseWriter, r *http.Request) *serverError {
	return writeError(w, errNotFound)
}

// Handles user creation
// POST /users
func createUserHandler(w http.ResponseWriter, r *http.Request) *serverError {
	var user model.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return &serverError{err, "can't decode req"}
	}

	if user.Email == nil || !(len(*user.Email) >= 3) || !(len(*user.Email) <= 256) {
		return writeError(w, errInputValidation)
	}

	if user.Password == nil || !(len(*user.Password) >= 6) || !(len(*user.Password) <= 258) {
		return writeError(w, errInputValidation)
	}

	if err := user.HashPassword(appConfig.BcryptCost); err != nil {
		return &serverError{err, "could not hash user pw"}
	}

	user.Password = nil

	tx, err := db.Begin()
	if err != nil {
		return &serverError{err, "could not begin tx"}
	}

	stmt, err := tx.Prepare(`
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING uuid, email, created_at
	`)
	if err != nil {
		return &serverError{err, "could not prepare tx"}
	}

	if err = stmt.QueryRow(user.Email, user.PasswordHash).Scan(&user.Uuid, &user.Email, &user.CreatedAt); err != nil {
		return &serverError{err, "could not insert"}
	}

	err = tx.Commit()

	return writeJson(w, user)
}

// Handles updating a user's information
func updateUserHandler(w http.ResponseWriter, r *http.Request) *serverError {
	return nil
}

// cancels. does not delete
func deleteOrderHandler(w http.ResponseWriter, r *http.Request) *serverError {
	order, ok := context.Get(r, reqOrder).(model.Order)
	if !ok {
		return &serverError{errors.New("this should not happen"), "this should not happen"}
	}

	if err := globalMatchingEngine.Cancel(&order); err != nil {
		return writeError(w, errTooBusy)
	}

	return nil
}

// Lists orders associated with a market
// GET /markets/{currencyPair}/orders
func listOrderHandler(w http.ResponseWriter, r *http.Request) *serverError {
	stmt, err := db.Prepare(`
		SELECT uuid, market_uuid, size, initial_size, price, side, status, created_at
		FROM orders
		WHERE market_uuid = $1 ORDER BY created_at DESC 
	`)
	if err != nil {
		return &serverError{err, "could not prepare stmt"}
	}

	mkt, ok := context.Get(r, reqMarket).(model.Market)
	if !ok {
		return &serverError{errors.New("this should not happen"), "this should not happen"}
	}

	rows, err := stmt.Query(mkt.Uuid)
	if err != nil {
		return &serverError{err, "could not query"}
	}
	defer rows.Close()

	var orders []*model.Order

	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.Uuid, &order.MarketUuid, &order.Size, &order.InitialSize, &order.Price, &order.Side, &order.Status, &order.CreatedAt)
		if err != nil {
			return &serverError{err, "error somewhere"}
		}

		orders = append(orders, &order)
	}

	return writeJson(w, orders)
}

// GET /user/accounts
func getAccountsHandler(w http.ResponseWriter, r *http.Request) *serverError {
	user, ok := context.Get(r, reqUser).(model.User)
	if !ok {
		return &serverError{errors.New("this should not happen"), "this should not happen"}
	}

	stmt, err := db.Prepare("SELECT uuid, user_uuid FROM accounts WHERE user_uuid = $1")
	if err != nil {
		return &serverError{err, "err preparing get accounts"}
	}

	var accounts []*model.Account

	rows, err := stmt.Query(user.Uuid)
	if err != nil {
		return &serverError{err, "err getting accts"}
	}
	defer rows.Close()

	for rows.Next() {
		var account model.Account
		err = rows.Scan(&account.Uuid, &account.UserUuid)
		if err != nil {
			return &serverError{err, "err getting acct"}
		}

		accounts = append(accounts, &account)
	}

	return writeJson(w, accounts)
}

// Handles getting the information of an order
// GET /accounts/{accountUuid}/orders/{orderUuid}
func getOrderHandler(w http.ResponseWriter, r *http.Request) *serverError {
	orderUuid := mux.Vars(r)["orderUuid"]

	stmt, err := db.Prepare(`
		SELECT uuid, market_uuid, size, initial_size, price, side, status, created_at
		FROM orders
		WHERE uuid = $1
	`)
	if err != nil {
		return &serverError{err, "could not prepare stmt"}
	}

	var order model.Order
	err = stmt.QueryRow(orderUuid).Scan(
		&order.Uuid,
		&order.Size,
		&order.InitialSize,
		&order.Price, order.Side,
		&order.Status,
		&order.CreatedAt)
	if err != nil {
		return &serverError{err, "could not get order values"}
	}

	return writeJson(w, order)
}

// Handles the creation of an order
// POST /accounts/{accountUuid}/orders
func createOrderHandler(w http.ResponseWriter, r *http.Request) *serverError {
	//market, ok := context.Get(r, reqMarket).(model.Market)
	//if !ok {
	//	return &serverError{errors.New("errors"), "error"}
	//}
	var order model.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		return &serverError{err, "could not decode input createOrderHandler"}
	}
	if order.Size == nil || !(*order.Size >= money.Satoshi) || !(*order.Size <= money.Bitcoin*1000) {
		return writeError(w, errInputValidation)
	}

	if order.Side == nil && (*order.Side != model.AskSide || *order.Side != model.BidSide) {
		return writeError(w, errInputValidation)
	}

	if order.Price == nil || !(*order.Price >= money.Satoshi) || !(*order.Price <= money.Bitcoin*1000) {
		return writeError(w, errInputValidation)
	}
	log.Println(order.MarketUuid)
	tx, err := db.Begin()
	if err != nil {
		return &serverError{err, "cannot begin tx"}
	}
	var market model.Market
	marketStmt, err := tx.Prepare(`SELECT uuid,base_currency,quote_currency,currency_pair FROM markets WHERE uuid = $1`)
	if err != nil {
		return &serverError{err, "euo"}
	}
	if err = marketStmt.QueryRow(order.MarketUuid).Scan(&market.Uuid, &market.BaseCurrency, &market.QuoteCurrency, &market.CurrencyPair); err != nil {
		return &serverError{err, "cannot begin tx"}
	}
	stmt, err := tx.Prepare(`
		INSERT INTO orders (market_uuid, size, initial_size, price, side, status, account_uuid)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING uuid, market_uuid, size, initial_size, price, side, status, created_at, account_uuid
	`)
	if err != nil {
		return &serverError{err, "error"}
	}

	err = stmt.QueryRow(
		market.Uuid,
		order.Size,
		order.Size,
		order.Price,
		order.Side,
		order.Status,
		context.Get(r, reqAccount).(model.Account).Uuid,
	).Scan(
		&order.Uuid,
		&order.MarketUuid,
		&order.Size,
		&order.InitialSize,
		&order.Price,
		&order.Side,
		&order.Status,
		&order.CreatedAt,
		&order.AccountUuid,
	)
	if err != nil {
		return &serverError{err, "error"}
	}
	stmt, err = tx.Prepare(`UPDATE balances SET available_balance = available_balance - $1, reserved_balance = reserved_balance + $1 WHERE currency = $2 AND account_uuid = $3 RETURNING available_balance`)
	if err != nil {
		return &serverError{err, "lol"}
	}
	if *order.Side == model.BidSide {
		var afterBalance *money.Unit
		//log.Println((order.InitialSize * *order.Price), market.QuoteCurrency, *order.AccountUuid)
		if err = stmt.QueryRow((order.InitialSize * *order.Price), market.QuoteCurrency, *order.AccountUuid).Scan(&afterBalance); err != nil {
			return &serverError{err, "111"}
		}
		if *afterBalance < money.Unit(0) {
			return writeError(w, errInsufficientFunds) // change to not enough moneyzzz
		}
	} else {
		var afterBalance *money.Unit
		//log.Println((order.InitialSize * *order.Price), market.baseCurrency, *order.AccountUuid)
		if err = stmt.QueryRow(order.InitialSize, market.BaseCurrency, *order.AccountUuid).Scan(&afterBalance); err != nil {
			return &serverError{err, "111"}
		}
		if *afterBalance < money.Unit(0) {
			return writeError(w, errInsufficientFunds) // change to not enough moneyzzz
		}
	}

	if err = tx.Commit(); err != nil {
		return &serverError{err, "tx commit err"}
	}
	globalMatchingEngine.Add(&order)

	return writeJson(w, order)
}
