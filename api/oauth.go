package main

import (
	"database/sql"
	"github.com/bitnel/api/model"
	"net/http"
	"github.com/gorilla/context"
	"time"
)

type oauthAccessToken struct {
	Uuid        string    `json:"-"`
	UserUuid    string    `json:"-"`
	AccessToken string    `json:"access_token"`
	ExpiresIn   int64     `json:"expires_in"`
	CreatedAt   time.Time `json:"-"`
	Scope 		string `json:"scope"`
}

func oauthTokenHandler(w http.ResponseWriter, r *http.Request) *serverError {
	if err := r.ParseForm(); err != nil {
		return &serverError{err, "unable to r.ParseForm()"}
	}
	if r.Form.Get("grant_type") != "password" {
		return writeError(w, errInputValidation)
	}
	var email, password string
	if email = r.Form.Get("username"); email == "" {
		return writeError(w, errInputValidation)
	}
	if password = r.Form.Get("password"); password == "" {
		return writeError(w, errInputValidation)
	}
	scope := r.Form.Get("scope");

	stmt, err := db.Prepare("SELECT uuid, password_hash FROM users WHERE email = $1")
	if err != nil {
		return &serverError{err, "unable to prepare stmt"}
	}

	var user model.User
	if err = stmt.QueryRow(email).Scan(&user.Uuid, &user.PasswordHash); err == sql.ErrNoRows {
		return writeError(w, errNotFound)
	} else if err != nil {
		return &serverError{err, "QueryRow error"}
	}

	if !user.ComparePassword(password) {
		return writeError(w, errNotFound)
	}

	stmt, err = db.Prepare(`INSERT INTO oauth_tokens (user_uuid, access_token, expires_at, created_at, scope)
		VALUES ($1, uuid_generate_v4(), NOW()+'1 day'::interval, NOW(), $2)
		RETURNING access_token, EXTRACT(EPOCH FROM (expires_at-NOW())), scope`)
	if err != nil {
		return &serverError{err, "oauth_token insert stmt prepare error"}
	}

	var token oauthAccessToken
	if err := stmt.QueryRow(user.Uuid, scope).Scan(&token.AccessToken, &token.ExpiresIn, &token.Scope); err != nil {
		return &serverError{err, "unable to exec oauth_tokens insert"}
	}
	context.Set(r,reqToken,token)

	return writeJson(w, token)
}
