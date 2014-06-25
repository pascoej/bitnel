package main

import (
	"database/sql"
	"errors"
	"github.com/bitnel/bitnel-api/model"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

// This type is used for keys for the context library
type contextVar int

const (
	reqMarket contextVar = iota
	reqUser
	reqAccount
	reqOrder
	reqToken
)

// This middleware wraps around all handlers concerning markets.
func marketFinder(fn apiHandler) apiHandler {
	return func(w http.ResponseWriter, r *http.Request) *serverError {
		pair := mux.Vars(r)["currencyPair"]

		stmt, err := db.Prepare(`SELECT uuid, base_currency, quote_currency, currency_pair
		FROM markets
		WHERE currency_pair = $1`)
		if err != nil {
			return &serverError{err, "could not prepare stmt"}
		}

		var market model.Market

		switch err = stmt.QueryRow(pair).Scan(&market.Uuid, &market.BaseCurrency, &market.QuoteCurrency, &market.CurrencyPair); {
		case err == sql.ErrNoRows:
			return writeError(w, errAuth)
		case err != nil:
			return &serverError{err, "could not get rows"}
		}

		context.Set(r, reqMarket, market)

		return fn(w, r)
	}
}

func accountFinder(fn apiHandler) apiHandler {
	return func(w http.ResponseWriter, r *http.Request) *serverError {
		uuid := mux.Vars(r)["accountUuid"]

		requestedUser, ok := context.Get(r, reqUser).(model.User)
		if !ok {
			return &serverError{errors.New("wtf happeend"), "wtf happened"}
		}

		stmt, err := db.Prepare(`SELECT uuid,user_uuid FROM accounts WHERE uuid = $1`)
		if err != nil {
			return &serverError{err, "err preparing acct getter"}
		}

		var account model.Account

		switch err = stmt.QueryRow(uuid).Scan(&account.Uuid, &account.UserUuid); {
		case err == sql.ErrNoRows:
			return writeError(w, errInputValidation)
		case err != nil:
			return &serverError{err, "err checking acct uuid"}
		}

		if account.UserUuid != requestedUser.Uuid {
			return writeError(w, errInputValidation)
		}

		context.Set(r, reqAccount, account)

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

		stmt, err := db.Prepare(`SELECT uuid, email, created_at
			FROM users WHERE uuid = (
			SELECT user_uuid FROM oauth_tokens WHERE access_token = $1 AND expires_at > NOW())`)
		if err != nil {
			return &serverError{err, "could not prepare stmt"}
		}

		var user model.User
		err = stmt.QueryRow(bearer[1]).Scan(&user.Uuid, &user.Email, &user.CreatedAt)

		switch {
		case err == sql.ErrNoRows:
			return writeError(w, errAuth)
		case err != nil:
			return &serverError{err, "could not get rows"}
		}

		context.Set(r, reqUser, user)

		return fn(w, r)
	}
}

func orderFinder(fn apiHandler) apiHandler {
	return func(w http.ResponseWriter, r *http.Request) *serverError {
		orderUuid := mux.Vars(r)["orderUuid"]

		requestedAccount, ok := context.Get(r, reqAccount).(model.Account)
		if !ok {
			return &serverError{errors.New("wtf happeend"), "wtf happened"}
		}

		stmt, err := db.Prepare(`SELECT uuid FROM orders WHERE uuid = $1 AND account_uuid = $2`)
		if err != nil {
			return &serverError{err, "err preparing order getter"}
		}

		var order model.Order

		switch err = stmt.QueryRow(orderUuid, requestedAccount.Uuid).Scan(&order.Uuid); {
		case err == sql.ErrNoRows:
			return writeError(w, errInputValidation)
		case err != nil:
			return &serverError{err, "err checking acct uuid"}
		}

		context.Set(r, reqOrder, order)

		return fn(w, r)
	}
}
