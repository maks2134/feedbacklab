// Package main is the entry point for the FeedbackLab API application.
package main

import (
	"innotech/internal/app"
	"innotech/internal/container"
	"innotech/pkg/db"
	"innotech/pkg/logger"
	"log"
)

// @title FeedbackLab API
// @version 0.1
// @host localhost:8080
// @BasePath /api
func main() {
	logger.Init()
	logger.Info("starting FeedbackLab API application")

	c := container.New()
	logger.Info("container initialized")

	if err := db.Migrate(c.Config.DatabaseURL, c.Config.MigrationsDir); err != nil {
		logger.Error("database migration failed",
			"error", err,
			"database_url", c.Config.DatabaseURL,
			"migrations_dir", c.Config.MigrationsDir,
		)
		log.Fatalf("Migration failed: %v", err)
	}
	logger.Info("database migrations completed successfully")

	logger.Info("starting application server",
		"version", "0.1",
	)

	app.Start(c)
}
