package app

import (
	_ "innotech/docs"
	"innotech/internal/contract"
	"innotech/internal/documentations"
	"innotech/internal/files"
	"innotech/internal/projects"
	"innotech/internal/user_projects"

	"innotech/internal/container"
	"innotech/internal/health"
	"innotech/internal/message_attachments"
	"innotech/internal/ticket_attachments"
	"innotech/internal/ticket_chats"
	"innotech/internal/tickets"

	"log"

	"github.com/gofiber/fiber/v2"
	swagger "github.com/swaggo/fiber-swagger"
)

func Start(container *container.Container) {
	app := fiber.New()

	app.Get("/swagger/*", swagger.WrapHandler)

	health.RegisterRoutes(app, container.HealthHandler)

	tickets.RegisterRoutes(app, container.TicketHandler)
	ticket_chats.RegisterRoutes(app, container.TicketChatsHandler)
	ticket_attachments.RegisterRoutes(app, container.TicketAttachmentsHandler)
	message_attachments.RegisterRoutes(app, container.MessageAttachmentsHandler)
	contract.RegisterRoutes(app, container.ContractHandler)
	projects.RegisterRoutes(app, container.ProjectHandler)
	documentations.RegisterRoutes(app, container.DocumentationHandler)
	user_projects.RegisterRoutes(app, container.UserProjectHandler)

	files.RegisterRoutes(app, container.FileHandler)

	log.Printf(" Server running on port %s\n", container.Config.AppPort)
	if err := app.Listen(":" + container.Config.AppPort); err != nil {
		log.Fatalf("failed to start feedbacklab: %v", err)
	}
}
