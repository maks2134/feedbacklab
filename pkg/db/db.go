package db

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func Connect(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", databaseURL)
	if err != nil {
		return nil, err
	}
	log.Println("database connected successfully")
	return db, nil
}

// изучить slog
