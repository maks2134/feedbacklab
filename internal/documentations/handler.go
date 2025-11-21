package documentations

import (
	"innotech/internal/storage/postgres"
	"innotech/internal/storage/transport"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for documentation operations.
type Handler struct {
	service Service
}

// NewHandler creates a new Handler instance.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create handles the creation of a new documentation entry.
func (h *Handler) Create(c *fiber.Ctx) error {
	dto := c.Locals("body").(*transport.CreateDocumentationDTO)
	d := postgres.Documentation{
		ProjectID:  dto.ProjectID,
		FilePath:   dto.FilePath,
		Version:    dto.Version,
		UploadedBy: dto.UploadedBy,
	}
	if err := h.service.Create(c.Context(), &d); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(d)
}

// GetByID retrieves a documentation entry by its ID.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	d, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(d)
}

// GetAll retrieves all documentation entries.
func (h *Handler) GetAll(c *fiber.Ctx) error {
	docs, err := h.service.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(docs)
}

// Update handles the update of an existing documentation entry.
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	dto := c.Locals("body").(*transport.UpdateDocumentationDTO)
	d := postgres.Documentation{
		ID:         id,
		FilePath:   dto.FilePath,
		Version:    dto.Version,
		UploadedBy: dto.UploadedBy,
	}
	if err := h.service.Update(c.Context(), &d); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(d)
}

// Delete handles the deletion of a documentation entry.
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	if err := h.service.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
