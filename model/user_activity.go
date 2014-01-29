package model

import (
	"net"
	"time"
)

type UserActivityType int

type UserActivity struct {
	Uuid      `json:"uuid"`
	UserUuid  `json:"user_uuid"`
	SourceIP  net.IP    `json:"source_ip"`
	CreatedAt time.Time `json:"created_at"`
}
