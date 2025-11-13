package handler

import (
	"innotech/internal/models"
	"innotech/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProjectHandler struct {
	service *service.ProjectService
}

func NewProjectHandler(service *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (h *ProjectHandler) RegisterRoutes(r fiber.Router) {
	group := r.Group("/projects")

	group.Get("/", h.GetAll)
	group.Get("/:id", h.GetByID)
	group.Post("/", h.Create)
	group.Put("/:id", h.Update)
	group.Delete("/:id", h.Delete)
}

// GetAll godoc
// @Summary Get all projects
// @Description Returns all projects
// @Tags Projects
// @Produce json
// @Success 200 {array} models.Project
// @Router /api/projects [get]
func (h *ProjectHandler) GetAll(c *fiber.Ctx) error {
	items, err := h.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(items)
}

// GetByID godoc
// @Summary Get project by ID
// @Description Returns a single project by its ID
// @Tags Projects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} models.Project
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/projects/{id} [get]
func (h *ProjectHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	item, err := h.service.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(item)
}

// Create godoc
// @Summary Create a new project
// @Description Creates a new project entry
// @Tags Projects
// @Accept json
// @Produce json
// @Param project body models.Project true "Project Data"
// @Success 201 {object} models.Project
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/projects [post]
func (h *ProjectHandler) Create(c *fiber.Ctx) error {
	var input models.Project
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.service.Create(&input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(input)
}

// Update godoc
// @Summary Update existing project
// @Description Updates project details by ID
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param project body models.Project true "Updated Project Data"
// @Success 200 {object} models.Project
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/projects/{id} [put]
func (h *ProjectHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var input models.Project
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	input.ID = id

	if err := h.service.Update(&input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(input)
}

// Delete godoc
// @Summary Delete project
// @Description Deletes a project by ID
// @Tags Projects
// @Param id path int true "Project ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/projects/{id} [delete]
func (h *ProjectHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	if err := h.service.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
