package main

import (
	"database/sql"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
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

func oauthTokenUserFinder(fn apiHandler) apiHandler {
	return func(w http.ResponseWriter, r *http.Request) *serverError {
		var authHeader string
		if authHeader = r.Header.Get("Authorization"); authHeader == "" {
			return writeError(w, errInputValidation)
		}

		var bearer []string
		if bearer = strings.Split(authHeader, " "); len(bearer) != 2 {
			return writeError(w, errInputValidation)
		}

		stmt, err := db.Prepare(`SELECT user_uuid FROM oauth_tokens WHERE access_token = $1`)
		if err != nil {
			return &serverError{err, "could not prepare stmt"}
		}

		log.Println(authHeader)

		var uuid string
		err = stmt.QueryRow(bearer[1]).Scan(&uuid)

		switch {
		case err == sql.ErrNoRows:
			log.Println("asdf")
			return writeError(w, errInputValidation)
		case err != nil:
			return &serverError{err, "could not get rows"}
		}

		context.Set(r, userUuid, uuid)
	}
}
