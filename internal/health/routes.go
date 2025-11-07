package health

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, h *Handler) {
	app.Get("/health", h.CheckHealth)
}
