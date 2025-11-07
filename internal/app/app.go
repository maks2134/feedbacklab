package app

import (
	"innotech/internal/container"
	"innotech/internal/health"
	"innotech/internal/message_attachments"
	"innotech/internal/ticket_attachments"
	"innotech/internal/ticket_chats"
	"innotech/internal/tickets"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Start(container *container.Container) {
	app := fiber.New()

	health.RegisterRoutes(app, container.HealthHandler)
	tickets.RegisterRoutes(app, container.TicketHandler)
	ticket_chats.RegisterRoutes(app, container.TicketChatsHandler)
	ticket_attachments.RegisterRoutes(app, container.TicketAttachmentsHandler)
	message_attachments.RegisterRoutes(app, container.MessageAttachmentsHandler)

	log.Printf("server running on port %s\n", container.Config.AppPort)
	if err := app.Listen(":" + container.Config.AppPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
