package container

import (
	"innotech/config"
	"innotech/internal/health"
	"innotech/internal/ticket_attachments"
	"innotech/internal/ticket_chats"
	"innotech/internal/tickets"
	"innotech/pkg/db"
	"log"

	"github.com/jmoiron/sqlx"
)

type Container struct {
	Config                   *config.Config
	DB                       *sqlx.DB
	HealthHandler            *health.Handler
	TicketHandler            *tickets.Handler
	TicketChatsHandler       *ticket_chats.Handler
	TicketAttachmentsHandler *ticket_attachments.Handler
}

func New() *Container {
	cfg := config.Load()

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
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

	return &Container{
		Config:                   cfg,
		DB:                       database,
		HealthHandler:            healthHandler,
		TicketHandler:            ticketHandler,
		TicketChatsHandler:       chatHandler,
		TicketAttachmentsHandler: attachHandler,
	}
}
