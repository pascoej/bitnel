package main

import (
	"database/sql"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
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
			return writeError(w, errInputValidation)
		case err != nil:
			return &serverError{err, "could not get rows"}
		}

		context.Set(r, marketUuid, uuid)

		return fn(w, r)
	}
}
