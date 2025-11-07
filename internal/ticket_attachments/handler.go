package ticket_attachments

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
// @Summary создать вложение тикета
// @Tags TicketAttachments
// @Accept json
// @Produce json
// @Param attachment body TicketAttachment true "Attachment"
// @Success 201 {object} TicketAttachment
// @Router /ticket_attachments/ [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var att TicketAttachment
	if err := c.BodyParser(&att); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if err := h.service.Create(&att); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(att)
}

// GetByID godoc
// @Summary получить вложение по ID
// @Tags TicketAttachments
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} TicketAttachment
// @Router /ticket_attachments/{id} [get]
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	att, err := h.service.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(att)
}

// GetByTicketID godoc
// @Summary получить все вложения тикета
// @Tags TicketAttachments
// @Produce json
// @Param ticket_id path int true "Ticket ID"
// @Success 200 {array} TicketAttachment
// @Router /ticket_attachments/ticket/{ticket_id} [get]
func (h *Handler) GetByTicketID(c *fiber.Ctx) error {
	ticketID, err := strconv.Atoi(c.Params("ticket_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ticket id"})
	}
	list, err := h.service.GetByTicketID(ticketID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}

// Update godoc
// @Summary обновить вложение
// @Tags TicketAttachments
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param attachment body TicketAttachment true "Attachment"
// @Success 200 {object} TicketAttachment
// @Router /ticket_attachments/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	var att TicketAttachment
	if err := c.BodyParser(&att); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	att.ID = id
	if err := h.service.Update(&att); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(att)
}

// Delete godoc
// @Summary удалить вложение
// @Tags TicketAttachments
// @Param id path int true "ID"
// @Success 204
// @Router /ticket_attachments/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	if err := h.service.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
