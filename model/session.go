package model

import (
	"time"
)

type Session struct {
	Uuid      string    `json:"-"`
	UserUuid  string    `json:"-"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
}
