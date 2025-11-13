package projects

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
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
	if err := h.service.Create(c.Context(), &p); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(p)
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	p, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(p)
}

func (h *Handler) GetAll(c *fiber.Ctx) error {
	ps, err := h.service.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(ps)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
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
	if err := h.service.Update(c.Context(), &p); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(p)
}

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
