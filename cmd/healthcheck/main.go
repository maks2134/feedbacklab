package main

import (
	"fmt"
	"innotech/config"
	"innotech/internal/health"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.Load()

	target := fmt.Sprintf("http://app:%s/api/health", cfg.AppPort)
	service := health.NewHTTPHealthService(target, 3*time.Second)
	handler := health.NewHandler(service)

	app := fiber.New()
	health.RegisterRoutes(app, handler)

	port := cfg.HealthPort

	log.Printf("healthcheck service running on port %s (target: %s)\n", port, target)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("failed to start healthcheck service: %v", err)
	}
}
