package user_projects

import (
	"innotech/internal/storage/transport"
	"innotech/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, h *Handler) {
	api := app.Group("/api/user-projects")

	api.Get("/", h.GetAll)
	api.Get("/:user_id/:project_id", h.Get)
	api.Post("/", middleware.ValidateBody[transport.CreateUserProjectDTO](h.Create))
	api.Put("/:user_id/:project_id", middleware.ValidateBody[transport.UpdateUserProjectDTO](h.Update))
	api.Delete("/:user_id/:project_id", h.Delete)
}
