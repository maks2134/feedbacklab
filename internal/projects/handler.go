package projects

import (
	"innotech/pkg/logger"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	logger.Info("project handler initialized")
	return &Handler{service: service}
}

func (h *Handler) Create(c *fiber.Ctx) error {
	dto := c.Locals("body").(*CreateProjectDTO)
	p := Project{
		Name:            dto.Name,
		Description:     dto.Description,
		GitlabProjectID: dto.GitlabProjectID,
		MattermostTeam:  dto.MattermostTeam,
	}

	logger.Info("handler: create project request",
		"name", p.Name,
		"path", c.Path(),
	)

	if err := h.service.Create(c.Context(), &p); err != nil {
		logger.Error("handler: create project failed",
			"error", err.Error(),
			"name", p.Name,
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Info("handler: project created",
		"id", p.ID,
		"name", p.Name,
	)
	return c.Status(fiber.StatusCreated).JSON(p)
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		logger.Warn("handler: invalid id param",
			"param", c.Params("id"),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	logger.Debug("handler: get project by id", "id", id)

	p, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		logger.Error("handler: get by id failed",
			"id", id,
			"error", err.Error(),
		)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Debug("handler: project retrieved successfully",
		"id", p.ID,
		"name", p.Name,
	)
	return c.JSON(p)
}

func (h *Handler) GetAll(c *fiber.Ctx) error {
	logger.Debug("handler: get all projects")

	ps, err := h.service.GetAll(c.Context())
	if err != nil {
		logger.Error("handler: get all failed",
			"error", err.Error(),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Info("handler: projects list retrieved",
		"count", len(ps),
	)
	return c.JSON(ps)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		logger.Warn("handler: invalid id param",
			"param", c.Params("id"),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	dto := c.Locals("body").(*UpdateProjectDTO)
	p := Project{
		ID:              id,
		Name:            dto.Name,
		Description:     dto.Description,
		GitlabProjectID: dto.GitlabProjectID,
		MattermostTeam:  dto.MattermostTeam,
	}

	logger.Info("handler: update project request",
		"id", p.ID,
		"name", p.Name,
	)

	if err := h.service.Update(c.Context(), &p); err != nil {
		logger.Error("handler: update project failed",
			"id", p.ID,
			"error", err.Error(),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Info("handler: project updated successfully",
		"id", p.ID,
		"name", p.Name,
	)
	return c.JSON(p)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		logger.Warn("handler: invalid id param",
			"param", c.Params("id"),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	logger.Warn("handler: delete project request", "id", id)

	if err := h.service.Delete(c.Context(), id); err != nil {
		logger.Error("handler: delete project failed",
			"id", id,
			"error", err.Error(),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Info("handler: project deleted successfully", "id", id)
	return c.SendStatus(fiber.StatusNoContent)
}
