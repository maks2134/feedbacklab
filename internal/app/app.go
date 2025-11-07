package app

import (
	"innotech/internal/container"
	"innotech/internal/health"
	"innotech/internal/tickets"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Start(container *container.Container) {
	app := fiber.New()

	health.RegisterRoutes(app, container.HealthHandler)
	tickets.RegisterRoutes(app, container.TicketHandler)

	log.Printf("server running on port %s\n", container.Config.AppPort)
	if err := app.Listen(":" + container.Config.AppPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
