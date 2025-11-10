package ticket_chats

import (
	"innotech/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, h *Handler) {
	api := app.Group("/api/ticket_chats")

	api.Get("/:id", h.GetByID)
	api.Get("/ticket/:ticket_id", h.GetByTicketID)
	api.Post("/", middleware.ValidateBody[CreateTicketChatDTO](h.Create))
	api.Put("/:id", middleware.ValidateBody[UpdateTicketChatDTO](h.Update))
	api.Delete("/:id", h.Delete)
}
