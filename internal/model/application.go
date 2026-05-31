package model

import "time"

type Application struct {
	ID           string     `db:"id"`
	SessionID    string     `db:"session_id"`
	ContactEmail string     `db:"contact_email"`
	StationName  string     `db:"station_name"`
	Description  string     `db:"description"`
	Genre        string     `db:"genre"`
	Notes        string     `db:"notes"`
	Status       string     `db:"status"`
	ReviewedBy   *string    `db:"reviewed_by"`
	ReviewedAt   *time.Time `db:"reviewed_at"`
	StationID    *string    `db:"station_id"`
	CreatedAt    time.Time  `db:"created_at"`
}
