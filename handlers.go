package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitnel/bitnel-api/model"
	"github.com/bitnel/bitnel-api/money"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
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

	err = stmt.QueryRow(*user.Email, user.PasswordHash).Scan(&user.Uuid, user.Email, &user.CreatedAt)
	if err != nil {
		return &serverError{err, "could not insert"}
	}

	return writeJson(w, user)
}

// Handles updating a user's information
func updateUserHandler(w http.ResponseWriter, r *http.Request) *serverError {
	return nil
}

// Handles getting the information of an order
func getOrderHandler(w http.ResponseWriter, r *http.Request) *serverError {
	orderUuid := mux.Vars(r)["orderUuid"]

	stmt, err := db.Prepare(`
		SELECT uuid, market_uuid, size, initial_size, price, side, status, type, created_at
		FROM orders
		WHERE uuid = $1
	`)
	if err != nil {
		return &serverError{err, "could not prepare stmt"}
	}

	var order model.Order
	err = stmt.QueryRow(orderUuid).Scan(&order.Uuid, order.Size, &order.InitialSize,
		order.Price, order.Side, order.Status, order.Type, &order.CreatedAt)
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

	if order.Size == nil || !(*order.Size >= money.Satoshi) || !(*order.Size <= money.Bitcoin*1000) {
		return writeError(w, errInputValidation)
	}

	if order.Side == nil && (*order.Side != model.AskSide || *order.Side != model.BidSide) {
		return writeError(w, errInputValidation)
	}

	if order.Type == nil && (*order.Type != model.MarketOrder || *order.Type != model.LimitOrder) {
		return writeError(w, errInputValidation)
	}

	tx, err := db.Begin()
	if err != nil {
		return &serverError{err, "cannot begin tx"}
	}

	switch *order.Type {
	case model.MarketOrder:
		order.Price = nil
	case model.LimitOrder:
		if order.Price == nil || !(*order.Price >= money.Satoshi) || !(*order.Size <= money.Bitcoin*1000) {
			return writeError(w, errInputValidation)
		}
	}

	stmt, err := tx.Prepare(`
		INSERT INTO orders (market_uuid, size, initial_size, price, side, status, type)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING uuid, market_uuid, size, initial_size, price, side, status, type, created_at
	`)
	if err != nil {
		return &serverError{err, "error"}
	}

	err = stmt.QueryRow(
		context.Get(r, marketUuid),
		*order.Size,
		*order.Size,
		*order.Price,
		*order.Side,
		order.Status,
		*order.Type,
	).Scan(
		&order.Uuid,
		&order.MarketUuid,
		order.Size,
		&order.InitialSize,
		order.Price,
		order.Side,
		&order.Status,
		order.Type,
		&order.CreatedAt,
	)

	if err != nil {
		return &serverError{err, "lol"}
	}

	return writeJson(w, order)
}
