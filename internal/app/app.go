package app

import (
	"innotech/config"
	"innotech/internal/handler"
	"innotech/internal/health"
	"innotech/internal/repository"
	"innotech/internal/service"
	"innotech/pkg/db"
	"log"

	"github.com/gofiber/fiber/v2"
)

func StartServer(port string) {
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

	app := fiber.New()

	health.RegisterRoutes(app)

	api := app.Group("/api")
	projectHandler.RegisterRoutes(api)
	moduleHandler.RegisterRoutes(api)
	contractHandler.RegisterRoutes(api)
	docHandler.RegisterRoutes(api)
	userProjectHandler.RegisterRoutes(api)

	log.Printf(" Server running on port %s\n", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
