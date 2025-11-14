package projects

import (
	"innotech/internal/storage/transport"
	"innotech/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, h *Handler) {
	api := app.Group("/api/projects")

	api.Get("/", h.GetAll)
	api.Get("/:id", h.GetByID)
	api.Post("/", middleware.ValidateBody[transport.CreateProjectDTO](h.Create))
	api.Put("/:id", middleware.ValidateBody[transport.UpdateProjectDTO](h.Update))
	api.Delete("/:id", h.Delete)
}
