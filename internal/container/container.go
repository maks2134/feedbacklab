// Package container provides dependency injection container for the application.
package container

import (
	"innotech/config"
	"innotech/internal/contract"
	"innotech/internal/documentations"
	"innotech/internal/files"
	"innotech/internal/health"
	"innotech/internal/messageattachments"
	"innotech/internal/modules"
	"innotech/internal/projects"
	"innotech/internal/ticketattachments"
	"innotech/internal/ticketchats"
	"innotech/internal/tickets"
	"innotech/internal/userprojects"
	"innotech/pkg/db"
	"innotech/pkg/i18n"
	"innotech/pkg/logger"
	"innotech/pkg/minio"
	"log"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

// Container holds all application dependencies and services.
type Container struct {
	Config                    *config.Config
	DB                        *sqlx.DB
	I18nBundle                *goi18n.Bundle
	HealthHandler             *health.Handler
	TicketHandler             *tickets.Handler
	TicketChatsHandler        *ticketchats.Handler
	TicketAttachmentsHandler  *ticketattachments.Handler
	MessageAttachmentsHandler *messageattachments.Handler
	ContractHandler           *contract.Handler
	ProjectHandler            *projects.Handler
	DocumentationHandler      *documentations.Handler
	Minio                     *minio.Client
	FileHandler               *files.Handler
	UserProjectHandler        *userprojects.Handler
	ModuleHandler             *modules.Handler
}

// New creates and initializes a new Container with all dependencies.
func New() *Container {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	newLogger := logger.NewLogger()
	slog.SetDefault(newLogger)
	// TODO сделать нормальный логгер через slog

	i18nBundle := i18n.InitBundle()
	localesDir := getEnv("LOCALES_DIR", "./locales")
	if err := i18n.LoadTranslations(i18nBundle, localesDir); err != nil {
		log.Printf("Warning: failed to load translations: %v", err)
	}

	healthService := health.NewSelfHealthService()
	healthHandler := health.NewHandler(healthService)

	ticketRepo := tickets.NewRepository(database)
	ticketService := tickets.NewService(ticketRepo)
	ticketHandler := tickets.NewHandler(ticketService)

	chatRepo := ticketchats.NewRepository(database)
	chatService := ticketchats.NewService(chatRepo)
	chatHandler := ticketchats.NewHandler(chatService)

	attachRepo := ticketattachments.NewRepository(database)
	attachService := ticketattachments.NewService(attachRepo)
	attachHandler := ticketattachments.NewHandler(attachService)

	msgAttachRepo := messageattachments.NewRepository(database)
	msgAttachService := messageattachments.NewService(msgAttachRepo)
	msgAttachHandler := messageattachments.NewHandler(msgAttachService)

	contractRepo := contract.NewRepository(database)
	contractService := contract.NewService(contractRepo)
	contractHandler := contract.NewHandler(contractService)

	projectRepo := projects.NewRepository(database)
	projectService := projects.NewService(projectRepo)
	projectHandler := projects.NewHandler(projectService)

	docRepo := documentations.NewRepository(database)
	docService := documentations.NewService(docRepo)
	docHandler := documentations.NewHandler(docService)

	userProjectRepo := userprojects.NewRepository(database)
	userProjectService := userprojects.NewService(userProjectRepo)
	userProjectHandler := userprojects.NewHandler(userProjectService)

	modulesRepo := modules.NewRepository(database)
	modulesService := modules.NewService(modulesRepo)
	modulesHandler := modules.NewHandler(modulesService)

	minioClient, err := minio.New(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.MinioBucket,
		cfg.MinioUseSSL,
	)

	fileService := files.NewService(minioClient, newLogger)
	fileHandler := files.NewHandler(fileService, newLogger)

	if err != nil {
		log.Fatalf("Failed to init MinIO: %v", err)
	}

	return &Container{
		Config:                    cfg,
		DB:                        database,
		I18nBundle:                i18nBundle,
		HealthHandler:             healthHandler,
		TicketHandler:             ticketHandler,
		TicketChatsHandler:        chatHandler,
		TicketAttachmentsHandler:  attachHandler,
		MessageAttachmentsHandler: msgAttachHandler,
		ContractHandler:           contractHandler,
		ProjectHandler:            projectHandler,
		DocumentationHandler:      docHandler,
		UserProjectHandler:        userProjectHandler,
		Minio:                     minioClient,
		FileHandler:               fileHandler,
		ModuleHandler:             modulesHandler,
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
