package handler

import (
	"innotech/internal/models"
	"innotech/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type DocumentationHandler struct {
	service *service.DocumentationService
}

func NewDocumentationHandler(service *service.DocumentationService) *DocumentationHandler {
	return &DocumentationHandler{service: service}
}

func (h *DocumentationHandler) RegisterRoutes(r fiber.Router) {
	group := r.Group("/documentations")

	group.Get("/", h.GetAll)
	group.Get("/:id", h.GetByID)
	group.Post("/", h.Create)
	group.Put("/:id", h.Update)
	group.Delete("/:id", h.Delete)
}

// GetAll godoc
// @Summary Get all documentations
// @Tags Documentations
// @Produce json
// @Success 200 {array} models.Documentation
// @Router /api/documentations [get]
func (h *DocumentationHandler) GetAll(c *fiber.Ctx) error {
	items, err := h.service.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(items)
}

// GetByID godoc
// @Summary Get documentation by ID
// @Tags Documentations
// @Produce json
// @Param id path int true "Documentation ID"
// @Success 200 {object} models.Documentation
// @Router /api/documentations/{id} [get]
func (h *DocumentationHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	item, err := h.service.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(item)
}

// Create godoc
// @Summary Create a new documentation
// @Tags Documentations
// @Accept json
// @Produce json
// @Param doc body models.Documentation true "Documentation Data"
// @Success 201 {object} models.Documentation
// @Router /api/documentations [post]
func (h *DocumentationHandler) Create(c *fiber.Ctx) error {
	var input models.Documentation
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.service.Create(&input); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(input)
}

// Update godoc
// @Summary Update existing documentation
// @Tags Documentations
// @Accept json
// @Produce json
// @Param id path int true "Documentation ID"
// @Param doc body models.Documentation true "Updated Documentation Data"
// @Success 200 {object} models.Documentation
// @Router /api/documentations/{id} [put]
func (h *DocumentationHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var input models.Documentation
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	input.ID = id
	if err := h.service.Update(&input); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(input)
}

// Delete godoc
// @Summary Delete documentation
// @Tags Documentations
// @Param id path int true "Documentation ID"
// @Success 204 "No Content"
// @Router /api/documentations/{id} [delete]
func (h *DocumentationHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	if err := h.service.Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
