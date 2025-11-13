package handler

import (
	"innotech/internal/models"
	"innotech/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ModuleHandler struct {
	service *service.ModuleService
}

func NewModuleHandler(service *service.ModuleService) *ModuleHandler {
	return &ModuleHandler{service: service}
}

func (h *ModuleHandler) RegisterRoutes(r fiber.Router) {
	group := r.Group("/modules")

	group.Get("/", h.GetAll)
	group.Get("/:id", h.GetByID)
	group.Post("/", h.Create)
	group.Put("/:id", h.Update)
	group.Delete("/:id", h.Delete)
}

// GetAll godoc
// @Summary Get all modules
// @Description Returns all modules
// @Tags Modules
// @Produce json
// @Success 200 {array} models.Module
// @Router /api/modules [get]
func (h *ModuleHandler) GetAll(c *fiber.Ctx) error {
	items, err := h.service.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(items)
}

// GetByID godoc
// @Summary Get module by ID
// @Description Returns a single module by its ID
// @Tags Modules
// @Produce json
// @Param id path int true "Module ID"
// @Success 200 {object} models.Module
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/modules/{id} [get]
func (h *ModuleHandler) GetByID(c *fiber.Ctx) error {
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
// @Summary Create a new module
// @Description Creates a new module entry
// @Tags Modules
// @Accept json
// @Produce json
// @Param module body models.Module true "Module Data"
// @Success 201 {object} models.Module
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/modules [post]
func (h *ModuleHandler) Create(c *fiber.Ctx) error {
	var input models.Module
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.service.Create(&input); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(input)
}

// Update godoc
// @Summary Update existing module
// @Description Updates module details by ID
// @Tags Modules
// @Accept json
// @Produce json
// @Param id path int true "Module ID"
// @Param module body models.Module true "Updated Module Data"
// @Success 200 {object} models.Module
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/modules/{id} [put]
func (h *ModuleHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var input models.Module
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
// @Summary Delete module
// @Description Deletes a module by ID
// @Tags Modules
// @Param id path int true "Module ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/modules/{id} [delete]
func (h *ModuleHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	if err := h.service.Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
