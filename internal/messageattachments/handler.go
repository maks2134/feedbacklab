// Package messageattachments provides message attachment management functionality.
package messageattachments

import (
	"innotech/internal/storage/postgres"
	"innotech/internal/storage/transport"
	"innotech/pkg/errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

// Handler handles HTTP requests for message attachment operations.
type Handler struct {
	service Service
}

// getLocalizer extracts localizer from context.
func getLocalizer(c *fiber.Ctx) *goi18n.Localizer {
	localizer, ok := c.Locals("localizer").(*goi18n.Localizer)
	if !ok || localizer == nil {
		// Fallback: return nil, will be handled in handleError
		return nil
	}
	return localizer
}

// NewHandler creates a new Handler instance.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary создать вложение сообщения
// @Tags MessageAttachments
// @Accept json
// @Produce json
// @Param attachment body transport.CreateMessageAttachmentDTO true "Attachment"
// @Success 201 {object} postgres.MessageAttachment
// @Router /message_attachments/ [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	localizer := getLocalizer(c)

	file, err := c.FormFile("file")
	if err != nil {
		msg, _ := localizer.Localize(&goi18n.LocalizeConfig{
			MessageID: "message_attachment.error.file_required",
		})
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
	}

	chatIDStr := c.FormValue("chat_id")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		msg, _ := localizer.Localize(&goi18n.LocalizeConfig{
			MessageID: "message_attachment.error.invalid_chat_id",
		})
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
	}

	uploadedByStr := c.FormValue("uploaded_by")
	uploadedBy, err := strconv.Atoi(uploadedByStr)
	if err != nil {
		msg, _ := localizer.Localize(&goi18n.LocalizeConfig{
			MessageID: "message_attachment.error.invalid_uploaded_by",
		})
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
	}

	att := postgres.MessageAttachment{
		ChatID:     chatID,
		UploadedBy: strconv.Itoa(uploadedBy),
	}

	if err := h.service.Create(c.Context(), &att, file); err != nil {
		return h.handleError(c, err, localizer)
	}

	return c.Status(fiber.StatusCreated).JSON(att)
}

// GetByID godoc
// @Summary получить вложение по ID
// @Tags MessageAttachments
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} postgres.MessageAttachment
// @Router /message_attachments/{id} [get]
func (h *Handler) GetByID(c *fiber.Ctx) error {
	localizer := getLocalizer(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		msg, _ := localizer.Localize(&goi18n.LocalizeConfig{
			MessageID: "common.error.bad_request",
		})
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
	}

	att, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return h.handleError(c, err, localizer)
	}
	return c.JSON(att)
}

// GetByChatID godoc
// @Summary получить все вложения по chat_id
// @Tags MessageAttachments
// @Produce json
// @Param chat_id path int true "Chat ID"
// @Success 200 {array} postgres.MessageAttachment
// @Router /message_attachments/chat/{chat_id} [get]
func (h *Handler) GetByChatID(c *fiber.Ctx) error {
	localizer := getLocalizer(c)

	chatID, err := strconv.Atoi(c.Params("chat_id"))
	if err != nil {
		msg, _ := localizer.Localize(&goi18n.LocalizeConfig{
			MessageID: "message_attachment.error.invalid_chat_id",
		})
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
	}

	list, err := h.service.GetByChatID(c.Context(), chatID)
	if err != nil {
		return h.handleError(c, err, localizer)
	}
	return c.JSON(list)
}

// Update godoc
// @Summary обновить вложение
// @Tags MessageAttachments
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param attachment body transport.UpdateMessageAttachmentDTO true "Attachment"
// @Success 200 {object} postgres.MessageAttachment
// @Router /message_attachments/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	localizer := getLocalizer(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		msg, _ := localizer.Localize(&goi18n.LocalizeConfig{
			MessageID: "common.error.bad_request",
		})
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
	}

	dto := c.Locals("body").(*transport.UpdateMessageAttachmentDTO)

	att := postgres.MessageAttachment{
		ID:       id,
		FilePath: dto.FilePath,
		FileType: dto.FileType,
	}

	if err := h.service.Update(c.Context(), &att); err != nil {
		return h.handleError(c, err, localizer)
	}
	return c.JSON(att)
}

// Delete godoc
// @Summary удалить вложение
// @Tags MessageAttachments
// @Param id path int true "ID"
// @Success 204
// @Router /message_attachments/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	localizer := getLocalizer(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		msg, _ := localizer.Localize(&goi18n.LocalizeConfig{
			MessageID: "common.error.bad_request",
		})
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return h.handleError(c, err, localizer)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// handleError processes application errors and returns appropriate HTTP response.
func (h *Handler) handleError(c *fiber.Ctx, err error, localizer *goi18n.Localizer) error {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		// Wrap unknown errors
		appErr = errors.NewInternalError("handler", "operation", err)
	}

	var statusCode int
	var msg string

	if localizer == nil {
		// Fallback to error message if localizer is not available
		msg = appErr.Message
	} else {
		switch appErr.Type {
		case errors.ErrorTypeNotFound:
			statusCode = fiber.StatusNotFound
			msg, _ = localizer.Localize(&goi18n.LocalizeConfig{
				MessageID: appErr.Key,
				DefaultMessage: &goi18n.Message{
					ID:    appErr.Key,
					Other: appErr.Message,
				},
			})
		case errors.ErrorTypeBadRequest, errors.ErrorTypeValidation:
			statusCode = fiber.StatusBadRequest
			msg, _ = localizer.Localize(&goi18n.LocalizeConfig{
				MessageID: appErr.Key,
				DefaultMessage: &goi18n.Message{
					ID:    appErr.Key,
					Other: appErr.Message,
				},
			})
		case errors.ErrorTypeUnauthorized:
			statusCode = fiber.StatusUnauthorized
			msg, _ = localizer.Localize(&goi18n.LocalizeConfig{
				MessageID: appErr.Key,
				DefaultMessage: &goi18n.Message{
					ID:    appErr.Key,
					Other: appErr.Message,
				},
			})
		case errors.ErrorTypeForbidden:
			statusCode = fiber.StatusForbidden
			msg, _ = localizer.Localize(&goi18n.LocalizeConfig{
				MessageID: appErr.Key,
				DefaultMessage: &goi18n.Message{
					ID:    appErr.Key,
					Other: appErr.Message,
				},
			})
		default:
			statusCode = fiber.StatusInternalServerError
			msg, _ = localizer.Localize(&goi18n.LocalizeConfig{
				MessageID: "common.error.internal",
			})
		}
	}

	if statusCode == 0 {
		// Set status code if not set
		switch appErr.Type {
		case errors.ErrorTypeNotFound:
			statusCode = fiber.StatusNotFound
		case errors.ErrorTypeBadRequest, errors.ErrorTypeValidation:
			statusCode = fiber.StatusBadRequest
		case errors.ErrorTypeUnauthorized:
			statusCode = fiber.StatusUnauthorized
		case errors.ErrorTypeForbidden:
			statusCode = fiber.StatusForbidden
		default:
			statusCode = fiber.StatusInternalServerError
		}
	}

	return c.Status(statusCode).JSON(fiber.Map{"error": msg})
}
