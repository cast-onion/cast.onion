package model

import "time"

type Station struct {
	ID          string    `db:"id"`
	Slug        string    `db:"slug"`
	DisplayName string    `db:"display_name"`
	Description *string   `db:"description"`
	Genre       *string   `db:"genre"`
	WebsiteURL  *string   `db:"website_url"`
	ArtKey      *string   `db:"art_key"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type User struct {
	ID           string    `db:"id"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	Role         string    `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
}
