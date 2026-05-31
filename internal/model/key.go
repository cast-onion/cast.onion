package model

import "time"

type StationKey struct {
	ID        string     `db:"id"`
	StationID string     `db:"station_id"`
	KeyHash   string     `db:"key_hash"`
	Revoked   bool       `db:"revoked"`
	CreatedAt time.Time  `db:"created_at"`
	RevokedAt *time.Time `db:"revoked_at"`
}
