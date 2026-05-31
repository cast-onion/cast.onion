package model

import "time"

type AdminAction struct {
	ID         string    `db:"id"`
	AdminID    string    `db:"admin_id"`
	Action     string    `db:"action"`
	TargetType string    `db:"target_type"`
	TargetID   string    `db:"target_id"`
	Reason     string    `db:"reason"`
	CreatedAt  time.Time `db:"created_at"`
}
