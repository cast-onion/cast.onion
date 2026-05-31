package main

import (
	"embed"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
)

//go:embed internal/db/migrations/*.sql
var migrationsFS embed.FS

func runMigrations(db *sqlx.DB) error {
	sub, err := fs.Sub(migrationsFS, "internal/db/migrations")
	if err != nil {
		return err
	}
	src, err := iofs.New(sub, ".")
	if err != nil {
		return err
	}
	driver, err := mysql.WithInstance(db.DB, &mysql.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", src, "mysql", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
