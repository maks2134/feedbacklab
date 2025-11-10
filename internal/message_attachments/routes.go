package message_attachments

import (
	"innotech/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, h *Handler) {
	api := app.Group("/api/message_attachments")

	api.Get("/:id", h.GetByID)
	api.Get("/chat/:chat_id", h.GetByChatID)
	api.Post("/", middleware.ValidateBody[CreateMessageAttachmentDTO](h.Create))
	api.Put("/:id", middleware.ValidateBody[UpdateMessageAttachmentDTO](h.Update))
	api.Delete("/:id", h.Delete)
}
