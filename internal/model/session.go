package model

import "time"

type Session struct {
	ID          string    `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	ExpiresAt   time.Time `db:"expires_at"`
	IPAddress   string    `db:"ip_address"`
	UserAgent   string    `db:"user_agent"`
	Invalidated bool      `db:"invalidated"`
}
