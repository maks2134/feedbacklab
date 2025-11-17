// Package user_projects provides user-project relationship management functionality.
package user_projects

import (
	"innotech/internal/storage/postgres"
	"innotech/internal/storage/transport"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for user-project operations.
type Handler struct {
	service Service
}

// NewHandler creates a new Handler instance.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create handles the creation of a new user-project relationship.
func (h *Handler) Create(c *fiber.Ctx) error {
	dto := c.Locals("body").(*transport.CreateUserProjectDTO)
	up := postgres.UserProject{
		UserID:      dto.UserID,
		ProjectID:   dto.ProjectID,
		Role:        dto.Role,
		Permissions: dto.Permissions,
	}
	if err := h.service.Create(c.Context(), &up); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(up)
}

// Get retrieves a user-project relationship by user ID and project ID.
func (h *Handler) Get(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	projectID, err := strconv.Atoi(c.Params("project_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid project id"})
	}
	up, err := h.service.Get(c.Context(), userID, projectID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(up)
}

// GetAll retrieves all user-project relationships.
func (h *Handler) GetAll(c *fiber.Ctx) error {
	list, err := h.service.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}

// Update handles the update of an existing user-project relationship.
func (h *Handler) Update(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	projectID, err := strconv.Atoi(c.Params("project_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid project id"})
	}
	dto := c.Locals("body").(*transport.UpdateUserProjectDTO)
	up := postgres.UserProject{
		UserID:      userID,
		ProjectID:   projectID,
		Role:        dto.Role,
		Permissions: dto.Permissions,
	}
	if err := h.service.Update(c.Context(), &up); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(up)
}

// Delete handles the deletion of a user-project relationship.
func (h *Handler) Delete(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	projectID, err := strconv.Atoi(c.Params("project_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid project id"})
	}
	if err := h.service.Delete(c.Context(), userID, projectID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
