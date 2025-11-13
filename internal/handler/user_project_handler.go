package handler

import (
	"innotech/internal/models"
	"innotech/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserProjectHandler struct {
	service *service.UserProjectService
}

func NewUserProjectHandler(service *service.UserProjectService) *UserProjectHandler {
	return &UserProjectHandler{service: service}
}

func (h *UserProjectHandler) RegisterRoutes(r fiber.Router) {
	group := r.Group("/user-projects")

	group.Get("/", h.GetAll)
	group.Get("/:id", h.GetByID)
	group.Post("/", h.Create)
	group.Put("/:id", h.Update)
	group.Delete("/:id", h.Delete)
}

// GetAll godoc
// @Summary Get all user-project relations
// @Tags UserProjects
// @Produce json
// @Success 200 {array} models.UserProject
// @Router /api/user-projects [get]
func (h *UserProjectHandler) GetAll(c *fiber.Ctx) error {
	items, err := h.service.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(items)
}

// GetByID godoc
// @Summary Get user-project by ID
// @Tags UserProjects
// @Produce json
// @Param id path int true "UserProject ID"
// @Success 200 {object} models.UserProject
// @Router /api/user-projects/{id} [get]
func (h *UserProjectHandler) GetByID(c *fiber.Ctx) error {
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
// @Summary Create a new user-project
// @Tags UserProjects
// @Accept json
// @Produce json
// @Param relation body models.UserProject true "UserProject Data"
// @Success 201 {object} models.UserProject
// @Router /api/user-projects [post]
func (h *UserProjectHandler) Create(c *fiber.Ctx) error {
	var input models.UserProject
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.service.Create(&input); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(input)
}

// Update godoc
// @Summary Update user-project relation
// @Tags UserProjects
// @Accept json
// @Produce json
// @Param id path int true "UserProject ID"
// @Param relation body models.UserProject true "Updated UserProject Data"
// @Success 200 {object} models.UserProject
// @Router /api/user-projects/{id} [put]
func (h *UserProjectHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var input models.UserProject
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
// @Summary Delete user-project
// @Tags UserProjects
// @Param id path int true "UserProject ID"
// @Success 204 "No Content"
// @Router /api/user-projects/{id} [delete]
func (h *UserProjectHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	if err := h.service.Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
