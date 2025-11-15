package container

import (
	"innotech/config"
	"innotech/internal/contract"
	"innotech/internal/documentations"
	"innotech/internal/health"
	"innotech/internal/message_attachments"
	"innotech/internal/projects"
	"innotech/internal/ticket_attachments"
	"innotech/internal/ticket_chats"
	"innotech/internal/tickets"
	"innotech/internal/user_projects"
	"innotech/pkg/db"
	"innotech/pkg/i18n"
	"innotech/pkg/logger"
	"log"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

type Container struct {
	Config                    *config.Config
	DB                        *sqlx.DB
	I18nBundle                *goi18n.Bundle
	HealthHandler             *health.Handler
	TicketHandler             *tickets.Handler
	TicketChatsHandler        *ticket_chats.Handler
	TicketAttachmentsHandler  *ticket_attachments.Handler
	MessageAttachmentsHandler *message_attachments.Handler
	ContractHandler           *contract.ContractHandler
	ProjectHandler            *projects.Handler
	DocumentationHandler      *documentations.Handler
	UserProjectHandler        *user_projects.Handler
}

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
	//TODO сделать нормальный логгер через slog

	// Инициализация i18n
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

	chatRepo := ticket_chats.NewRepository(database)
	chatService := ticket_chats.NewService(chatRepo)
	chatHandler := ticket_chats.NewHandler(chatService)

	attachRepo := ticket_attachments.NewRepository(database)
	attachService := ticket_attachments.NewService(attachRepo)
	attachHandler := ticket_attachments.NewHandler(attachService)

	msgAttachRepo := message_attachments.NewRepository(database)
	msgAttachService := message_attachments.NewService(msgAttachRepo)
	msgAttachHandler := message_attachments.NewHandler(msgAttachService)

	contractRepo := contract.NewContractRepository(database)
	contractService := contract.NewContractService(contractRepo)
	contractHandler := contract.NewContractHandler(contractService)

	projectRepo := projects.NewRepository(database)
	projectService := projects.NewService(projectRepo)
	projectHandler := projects.NewHandler(projectService)

	docRepo := documentations.NewRepository(database)
	docService := documentations.NewService(docRepo)
	docHandler := documentations.NewHandler(docService)

	userProjectRepo := user_projects.NewRepository(database)
	userProjectService := user_projects.NewService(userProjectRepo)
	userProjectHandler := user_projects.NewHandler(userProjectService)

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
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
