package app

import (
	"innotech/internal/container"
	"innotech/internal/health"
	"innotech/internal/message_attachments"
	"innotech/internal/ticket_attachments"
	"innotech/internal/ticket_chats"
	"innotech/internal/tickets"

	"innotech/config"
	"innotech/internal/handler"
	"innotech/internal/repository"
	"innotech/internal/service"
	"innotech/pkg/db"

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

	cfg := config.Load()

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer database.Close()

	projectRepo := repository.NewProjectRepository(database)
	moduleRepo := repository.NewModuleRepository(database)
	contractRepo := repository.NewContractRepository(database)
	docRepo := repository.NewDocumentationRepository(database)
	userProjectRepo := repository.NewUserProjectRepository(database)

	projectService := service.NewProjectService(projectRepo)
	moduleService := service.NewModuleService(moduleRepo)
	contractService := service.NewContractService(contractRepo)
	docService := service.NewDocumentationService(docRepo)
	userProjectService := service.NewUserProjectService(userProjectRepo)

	projectHandler := handler.NewProjectHandler(projectService)
	moduleHandler := handler.NewModuleHandler(moduleService)
	contractHandler := handler.NewContractHandler(contractService)
	docHandler := handler.NewDocumentationHandler(docService)
	userProjectHandler := handler.NewUserProjectHandler(userProjectService)

	api := app.Group("/api")
	projectHandler.RegisterRoutes(api)
	moduleHandler.RegisterRoutes(api)
	contractHandler.RegisterRoutes(api)
	docHandler.RegisterRoutes(api)
	userProjectHandler.RegisterRoutes(api)

	log.Printf(" Server running on port %s\n", container.Config.AppPort)
	if err := app.Listen(":" + container.Config.AppPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
