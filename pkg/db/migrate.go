package db

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Migrate(dbURL, migrationsDir string) error {
	db, err := goose.OpenDBWithDriver("pgx", dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}

	log.Println("migrations applied successfully")
	return nil
}
