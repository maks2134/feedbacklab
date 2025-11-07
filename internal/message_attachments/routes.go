package message_attachments

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, h *Handler) {
	r := app.Group("/api/message_attachments")
	r.Post("/", h.Create)
	r.Get("/:id", h.GetByID)
	r.Get("/chat/:chat_id", h.GetByChatID)
	r.Put("/:id", h.Update)
	r.Delete("/:id", h.Delete)
}
