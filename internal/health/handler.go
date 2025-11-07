package health

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service HealthService
}

func NewHandler(service HealthService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CheckHealth(c *fiber.Ctx) error {
	ctx := context.Background()
	if err := h.service.Check(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "unhealthy",
			"error":  err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "healthy",
	})
}
