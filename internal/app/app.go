package app

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func StartServer(port string) {
	app := fiber.New()

	//health.RegisterRoutes(app)

	log.Printf("server running on port %s\n", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
