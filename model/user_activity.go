package model

import (
	"net"
	"time"
)

type UserActivityType int

type UserActivity struct {
	Uuid      string    `json:"uuid"`
	UserUuid  string    `json:"user_uuid"`
	SourceIP  net.IP    `json:"source_ip"`
	CreatedAt time.Time `json:"created_at"`
}
