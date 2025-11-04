package main

import (
	"innotech/config"
	"innotech/internal/app"
	"innotech/pkg/db"
	"log"

	"github.com/jmoiron/sqlx"
)

func main() {
	cfg := config.Load()

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer func(database *sqlx.DB) {
		err := database.Close()
		if err != nil {
			log.Fatalf("DB connection failed: %v", err)
		}
	}(database)

	if err := db.Migrate(cfg.DatabaseURL, cfg.MigrationsDir); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	app.StartServer(cfg.AppPort)
}
