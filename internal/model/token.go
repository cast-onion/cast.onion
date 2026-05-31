package model

import "time"

type AccessToken struct {
	ID        string     `db:"id"`
	StationID string     `db:"station_id"`
	TokenHash string     `db:"token_hash"`
	Revoked   bool       `db:"revoked"`
	CreatedAt time.Time  `db:"created_at"`
	RevokedAt *time.Time `db:"revoked_at"`
}
