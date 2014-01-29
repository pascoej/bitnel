package main

import (
	"database/sql"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// This type is used for keys for the context library
type contextVar int

const (
	marketUuid contextVar = iota
	userUuid
)

// This middleware wraps around all handlers concerning markets.
func marketFinder(fn apiHandler) apiHandler {
	return func(w http.ResponseWriter, r *http.Request) *serverError {
		pair := mux.Vars(r)["currencyPair"]

		stmt, err := db.Prepare(`SELECT uuid FROM markets WHERE currency_pair = $1`)
		if err != nil {
			return &serverError{err, "could not prepare stmt"}
		}

		var uuid string

		err = stmt.QueryRow(pair).Scan(&uuid)

		switch {
		case err == sql.ErrNoRows:
			return writeError(w, errAuth)
		case err != nil:
			return &serverError{err, "could not get rows"}
		}

		context.Set(r, marketUuid, uuid)

		return fn(w, r)
	}
}
func sessionFinder(fn apiHandler) apiHandler {
	return func(w http.ResponseWriter, r *http.Request) *serverError {
		token := r.Header.Get("token")

		log.Println(token)
		stmt, err := db.Prepare(`SELECT user_uuid FROM sessions WHERE token = $1 AND expires_at > NOW()`)
		if err != nil {
			return &serverError{err, "err db"}
		}
		var uuid string
		err = stmt.QueryRow(token).Scan(&uuid)
		if err != nil {
			return writeError(w, errAuth)
		}
		context.Set(r, userUuid, uuid)
		return fn(w, r)
	}
}
