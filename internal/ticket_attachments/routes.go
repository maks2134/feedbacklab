package ticket_attachments

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, h *Handler) {
	r := app.Group("/api/ticket_attachments")
	r.Post("/", h.Create)
	r.Get("/:id", h.GetByID)
	r.Get("/ticket/:ticket_id", h.GetByTicketID)
	r.Put("/:id", h.Update)
	r.Delete("/:id", h.Delete)
}
