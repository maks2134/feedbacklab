package container

import (
	"innotech/config"
	"innotech/internal/health"
	"innotech/internal/tickets"
	"innotech/pkg/db"
	"log"

	"github.com/jmoiron/sqlx"
)

type Container struct {
	Config        *config.Config
	DB            *sqlx.DB
	HealthHandler *health.Handler
	TicketHandler *tickets.Handler
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

	return &Container{
		Config:        cfg,
		DB:            database,
		HealthHandler: healthHandler,
		TicketHandler: ticketHandler,
	}
}
