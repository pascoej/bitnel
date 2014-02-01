package main

import (
	"database/sql"
	"github.com/bitnel/bitnel-api/model"
	"net/http"
	"time"
)

type oauthAccessToken struct {
	Uuid        string    `json:"-"`
	UserUuid    string    `json:"-"`
	AccessToken string    `json:"access_token"`
	ExpiresAt   int64     `json:"expires_at"`
	CreatedAt   time.Time `json:"-"`
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

	stmt, err = db.Prepare(`INSERT INTO oauth_tokens (user_uuid, access_token, expires_at)
		VALUES ($1, uuid_generate_v4(), 3600
		RETURNING access_token, expires_at`)
	if err != nil {
		return &serverError{err, "oauth_token insert stmt prepare error"}
	}

	var token oauthAccessToken
	if err := stmt.QueryRow(user.Uuid).Scan(&token.AccessToken, &token.ExpiresAt); err != nil {
		return &serverError{err, "unable to exec oauth_tokens insert"}
	}

	return writeJson(w, token)
}
