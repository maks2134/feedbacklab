package contract

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, h *ContractHandler) {
	api := app.Group("/api/contracts")

	api.Get("/", h.GetAll)
	api.Get("/:id", h.GetByID)
	api.Post("/", h.Create)
	api.Put("/:id", h.Update)
	api.Delete("/:id", h.Delete)
}
