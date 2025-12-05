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
// @description API для управления проектами, тикетами и документацией
// @host localhost:8080
// @BasePath /api
// @schemes http
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

// todo 1) permission - написать библу, получение токена, распаршиваю, отправляю в кейлкоак, сделать метод, фиксирующий название, просто эндпоинт, тикет лист, получаю лист чекаю токен. 2) MinIO - подключить полностью к аттачментам, полный URL - адресс и подпись(бенчмарк - безопасная ссылка). 3) посмотреть интеграции с mattermost, поднять докер, и почекать докуму, и планировщик https://temporal.io/
