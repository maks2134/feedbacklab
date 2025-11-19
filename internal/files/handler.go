package files

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
	logger  *slog.Logger
}

func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// Upload POST /files/upload
func (h *Handler) Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		h.logger.Error("no file provided", "err", err)
		return fiber.NewError(fiber.StatusBadRequest, "file is required")
	}

	savePath := "/tmp/" + file.Filename

	if err := c.SaveFile(file, savePath); err != nil {
		h.logger.Error("cannot save uploaded file", "err", err)
		return fiber.NewError(fiber.StatusInternalServerError, "file save error")
	}

	url, err := h.service.UploadFile(c.Context(), file.Filename, savePath)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "upload to storage failed")
	}

	return c.JSON(fiber.Map{
		"url": url,
	})
}
