package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitnel/bitnel-api/model"
	"github.com/bitnel/bitnel-api/money"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
<<<<<<< HEAD
=======
	"time"
	//"github.com/gorilla/websocket"
>>>>>>> 8822ead2d45d8caa6d290ab78fc0e24a8ef488d4
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
	return nil
}

// Lists orders associated with a market
func listOrderHandler(w http.ResponseWriter, r *http.Request) *serverError {
	stmt, err := db.Prepare(`
		SELECT uuid, market_uuid, size, initial_size, price, side, status, created_at
		FROM orders
		WHERE market_uuid = $1 ORDER BY created_at DESC 
	`)
	if err != nil {
		return &serverError{err, "could not prepare stmt"}
	}

	rows, err := stmt.Query(context.Get(r, marketUuid))
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
func getAccountsHandler(w http.ResponseWriter, r *http.Request) *serverError {
	uUuid := context.Get(r, userUuid)
	stmt, err := db.Prepare(`SELECT uuid,user_uuid,type FROM accounts WHERE user_uuid = $1`)
	if err != nil {
		return &serverError{err, "err preparing get accounts"}
	}
	var accounts []model.Account
	rows, err := stmt.Query(uUuid)
	if err != nil {
		return &serverError{err, "err getting accts"}
	}
	for rows.Next() {
		var account model.Account
		err = rows.Scan(&account.Uuid, &account.UserUuid, &account.Type)
		if err != nil {
			return &serverError{err, "err getting acct"}
		}
		accounts = append(accounts, account)
	}
	return writeJson(w, accounts)
}

// Handles getting the information of an order
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
func createOrderHandler(w http.ResponseWriter, r *http.Request) *serverError {
	var order model.Order

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		return &serverError{err, "could not decode input createOrderHandler"}
	}
<<<<<<< HEAD
=======
	return writeJson(w, order)
}
func createOrder(marketUuid interface{}, order *model.Order) (*serverError, *model.Order) {
	return nil, order
}

// Handles the creation of an order
func createOrderHandler(w http.ResponseWriter, r *http.Request) *serverError {
	var order *model.Order

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		return &serverError{err, "could not decode input createOrderHandler"}
	}
>>>>>>> 8822ead2d45d8caa6d290ab78fc0e24a8ef488d4

	if order.Size == nil || !(*order.Size >= money.Satoshi) || !(*order.Size <= money.Bitcoin*1000) {
		return writeError(w, errInputValidation)
	}

	if order.Side == nil && (*order.Side != model.AskSide || *order.Side != model.BidSide) {
		return writeError(w, errInputValidation)
	}

	if order.Price == nil || !(*order.Price >= money.Satoshi) || !(*order.Price <= money.Bitcoin*1000) {
		return writeError(w, errInputValidation)
	}
<<<<<<< HEAD

=======
>>>>>>> 8822ead2d45d8caa6d290ab78fc0e24a8ef488d4
	tx, err := db.Begin()
	if err != nil {
		return &serverError{err, "cannot begin tx"}
	}
	stmt, err := tx.Prepare(`
		INSERT INTO orders (market_uuid, size, initial_size, price, side, status, account_uuid)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING uuid, market_uuid, size, initial_size, price, side, status, created_at,account_uuid
	`)
	if err != nil {
		return &serverError{err, "error"}
	}

	err = stmt.QueryRow(
		context.Get(r, marketUuid),
		order.Size,
		order.Size,
		order.Price,
		order.Side,
		order.Status,
		context.Get(r, userUuid),
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
<<<<<<< HEAD
		return &serverError{err, "lol"}
=======
		return &serverError{err, "err scanning columns"}
>>>>>>> 8822ead2d45d8caa6d290ab78fc0e24a8ef488d4
	}

	if err = tx.Commit(); err != nil {
		return &serverError{err, "tx commit err"}
	}
<<<<<<< HEAD

	globalMatchingEngine.Add(&order)

	return writeJson(w, order)
=======
	globalMatchingEngine.Add(order)
	return writeJson(w, order)
}

func createSessionHandler(w http.ResponseWriter, r *http.Request) *serverError {
	var user model.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return &serverError{err, "could not decode input createOrderHandler"}
	}

	if user.Email == nil {
		return writeError(w, errInputValidation)
	}

	if user.Password == nil {
		return writeError(w, errInputValidation)
	}
	stmt, err := db.Prepare(`SELECT password_hash FROM users WHERE email = $1`)
	if err != nil {
		return &serverError{err, "dfd"}
	}

	if err = stmt.QueryRow(user.Email).Scan(&user.PasswordHash); err != nil {
		return &serverError{err, "no such user"}
	}
	if !user.ComparePassword(*user.Password) {
		return writeError(w, errNotFound)
	}
	stmt, err = db.Prepare(`INSERT INTO sessions (user_uuid, expires_at) VALUES((SELECT uuid FROM users WHERE email = $1),$2) RETURNING uuid,user_uuid,token,created_at,expires_at`)
	if err != nil {
		return &serverError{err, "err preparing"}
	}
	var session model.Session
	if err := stmt.QueryRow(user.Email, time.Now().Add(time.Hour*24*3)).Scan(&session.Uuid, &session.UserUuid, &session.Token, &session.CreatedAt, &session.ExpiresAt); err != nil {
		return &serverError{err, "err creating token"}
	}

	return writeJson(w, session)
>>>>>>> 8822ead2d45d8caa6d290ab78fc0e24a8ef488d4
}
