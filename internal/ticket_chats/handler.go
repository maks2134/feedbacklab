package ticket_chats

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
	var chat TicketChat
	if err := c.BodyParser(&chat); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if err := h.service.Create(&chat); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(chat)
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	chat, err := h.service.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(chat)
}

func (h *Handler) GetByTicketID(c *fiber.Ctx) error {
	ticketID, err := strconv.Atoi(c.Params("ticket_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ticket id"})
	}
	chats, err := h.service.GetByTicketID(ticketID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(chats)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	var chat TicketChat
	if err := c.BodyParser(&chat); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	chat.ID = id
	if err := h.service.Update(&chat); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(chat)
}

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
