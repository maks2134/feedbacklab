package tickets

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, h *Handler) {
	api := app.Group("/api/tickets")
	api.Post("/", h.Create)
	api.Get("/", h.GetAll)
	api.Get("/:id", h.GetByID)
	api.Put("/:id", h.Update)
	api.Delete("/:id", h.Delete)
}
