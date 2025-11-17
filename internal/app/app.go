// Package app provides application initialization and routing.
package app

import (
	// Import swagger docs for API documentation.
	_ "innotech/docs"
	"innotech/internal/contract"
	"innotech/internal/modules"

	"innotech/internal/documentations"
	"innotech/internal/projects"
	"innotech/internal/userprojects"
	"strconv"

	"innotech/internal/container"
	"innotech/internal/health"
	"innotech/internal/messageattachments"
	"innotech/internal/ticketattachments"
	"innotech/internal/ticketchats"
	"innotech/internal/tickets"
	"innotech/pkg/middleware"

	"log"

	"github.com/gofiber/fiber/v2"
	swagger "github.com/swaggo/fiber-swagger"
)

// Start initializes and starts the HTTP server with all registered routes.
func Start(container *container.Container) {
	app := fiber.New()

	app.Use(middleware.I18nMiddleware(container.I18nBundle))

	app.Get("/swagger/*", swagger.WrapHandler)

	health.RegisterRoutes(app, container.HealthHandler)

	tickets.RegisterRoutes(app, container.TicketHandler)
	ticketchats.RegisterRoutes(app, container.TicketChatsHandler)
	ticketattachments.RegisterRoutes(app, container.TicketAttachmentsHandler)
	messageattachments.RegisterRoutes(app, container.MessageAttachmentsHandler)
	contract.RegisterRoutes(app, container.ContractHandler)
	projects.RegisterRoutes(app, container.ProjectHandler)
	documentations.RegisterRoutes(app, container.DocumentationHandler)
	userprojects.RegisterRoutes(app, container.UserProjectHandler)
	modules.RegisterRoutes(app, container.ModuleHandler)

	log.Printf(" Server running on port %d\n", container.Config.AppPort)
	if err := app.Listen(":" + strconv.Itoa(container.Config.AppPort)); err != nil {
		log.Fatalf("failed to start feedbacklab: %v", err)
	}
}
