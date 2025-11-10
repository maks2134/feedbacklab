package tickets

import (
	"innotech/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, h *Handler) {
	api := app.Group("/api/tickets")

	api.Get("/", h.GetAll)
	api.Get("/:id", h.GetByID)

	api.Post("/", middleware.ValidateBody[CreateTicketDTO](h.Create))
	api.Put("/:id", middleware.ValidateBody[UpdateTicketDTO](h.Update))

	api.Delete("/:id", h.Delete)
}
