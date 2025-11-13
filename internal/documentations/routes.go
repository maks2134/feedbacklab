package documentations

import (
	"innotech/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, h *Handler) {
	api := app.Group("/api/documentations")

	api.Get("/", h.GetAll)
	api.Get("/:id", h.GetByID)
	api.Post("/", middleware.ValidateBody[CreateDocumentationDTO](h.Create))
	api.Put("/:id", middleware.ValidateBody[UpdateDocumentationDTO](h.Update))
	api.Delete("/:id", h.Delete)
}
