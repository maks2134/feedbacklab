package container

import (
	"innotech/config"
	"innotech/internal/contract"
	"innotech/internal/health"
	"innotech/internal/message_attachments"
	"innotech/internal/ticket_attachments"
	"innotech/internal/ticket_chats"
	"innotech/internal/tickets"
	"innotech/pkg/db"
	"innotech/pkg/logger"
	"log"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Container struct {
	Config                    *config.Config
	DB                        *sqlx.DB
	HealthHandler             *health.Handler
	TicketHandler             *tickets.Handler
	TicketChatsHandler        *ticket_chats.Handler
	TicketAttachmentsHandler  *ticket_attachments.Handler
	MessageAttachmentsHandler *message_attachments.Handler
	ContractHandler           *contract.ContractHandler
	//DocumentationHandler	  *documentation.DocumentationHandler
}

func New() *Container {
	cfg := config.Load()

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	newLogger := logger.NewLogger()
	slog.SetDefault(newLogger)
	//TODO сделать нормальный логгер через slog

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

	//documRepo := documentation.NewDocumentationRepository(database)
	//documService := documentation.NewDocumentationService(documRepo)
	//documHandler : =documentation.NewDocumentationHandler(documService)

	return &Container{
		Config:                    cfg,
		DB:                        database,
		HealthHandler:             healthHandler,
		TicketHandler:             ticketHandler,
		TicketChatsHandler:        chatHandler,
		TicketAttachmentsHandler:  attachHandler,
		MessageAttachmentsHandler: msgAttachHandler,
		ContractHandler:           contractHandler,
		//DocumentationHandler:      documHandler,
	}
}
