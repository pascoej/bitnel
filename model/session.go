package model

import (
	"time"
	""
)

type Session struct {
	Uuid      string
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}
