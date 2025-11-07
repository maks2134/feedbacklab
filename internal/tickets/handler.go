package tickets

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

// Create godoc
// @Summary создать тикет
// @Tags Tickets
// @Accept json
// @Produce json
// @Param ticket body Ticket true "Ticket"
// @Success 201 {object} Ticket
// @Failure 400 {object} map[string]string
// @Router /tickets/ [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var t Ticket
	if err := c.BodyParser(&t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.Create(c.Context(), &t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(t)
}

// GetByID godoc
// @Summary получить тикет по ID
// @Tags Tickets
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} Ticket
// @Failure 404 {object} map[string]string
// @Router /tickets/{id} [get]
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	t, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(t)
}

// GetAll godoc
// @Summary получить все тикеты
// @Tags Tickets
// @Produce json
// @Success 200 {object} Ticket
// @Failure 404 {object} map[string]string
// @Router /tickets/ [get]
func (h *Handler) GetAll(c *fiber.Ctx) error {
	tickets, err := h.service.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(tickets)
}

// Update godoc
// @Summary обновить тикет
// @Tags Tickets
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param ticket body Ticket true "Ticket"
// @Success 200 {object} Ticket
// @Router /tickets/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	var t Ticket
	if err := c.BodyParser(&t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	t.ID = id

	if err := h.service.Update(c.Context(), &t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(t)
}

// Delete godoc
// @Summary удалить тикет
// @Tags Tickets
// @Param id path int true "ID"
// @Success 204
// @Router /tickets/{id} [delete]
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
