package main

import (
	"innotech/internal/app"
	"innotech/internal/container"
	"innotech/pkg/db"
	"log"
)

// @title FeedbackLab API
// @version 0.1
// @host localhost:8080
// @BasePath /api
func main() {
	c := container.New()

	if err := db.Migrate(c.Config.DatabaseURL, c.Config.MigrationsDir); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	app.Start(c)
}
