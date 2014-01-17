package main

import (
	//"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	"net/http"
)

// This type is used for keys for the context library
type contextVar int

const (
	marketUuid contextVar = iota
)

// This serves the root path of our API. Be friendly; say hello.
func indexHandler(w http.ResponseWriter, r *http.Request) *serverError {
	fmt.Fprintln(w, "Welcome to the Bitnel API!")

	return nil
}

// We define our own not found handler because we dislike the default Gorilla
// 404 message.
func notFoundHandler(w http.ResponseWriter, r *http.Request) *serverError {
	writeError(w, errNotFound)

	return nil
}

func createUserHandler(w http.ResponseWriter, r *http.Request) *serverError {
	return nil
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) *serverError {
	return nil
}

func createOrderHandler(w http.ResponseWriter, r *http.Request) *serverError {
	var order Order

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		return &serverError{err, "could not decode input createOrderHandle"}
	}

	if order.Size == nil {
		writeError(w, errInputValidation)
		return nil
	}

	writeJson(w, apiResponse{"uuid": context.Get(r, marketUuid)})

	/*

		tx, err := db.Begin()
		if err != nil {
			//handle error
		}

		stmt, err := stmt.Prepare(`
			INSERT INTO orders (market_uuid, size, initial_size, price, side, status, type)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING uuid, market_uuid, size, initial_size, price, side, status, type, created_at
		`)
		if err != nil {
			// handle error
		}

		var ao Order

		err = stmt.QueryRow(context.Get(r, marketUuid), order.Size, order.InitialSize, order.Price, order.Side, order.Status, order.Type).Scan(&ao.Uuid, &ao.MarketUuid,
	*/

	return nil
}
