package auth

import (
	"log/slog"
	"strings"

	"innotech/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

//TODO: перевод сваггера

// Handler handles authentication-related HTTP requests
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new auth handler
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// CreateUser handles user creation (admin only)
// @Summary Create a new user
// @Description Create a new user in Keycloak. Only administrators can perform this action.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User creation request"
// @Success 201 {object} CreateUserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/users [post]
func (h *Handler) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("failed to parse request body",
			"error", err,
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := middleware.ValidateBody[CreateUserRequest](func(c *fiber.Ctx) error {
		return nil
	})(c); err != nil {
		return err
	}

	body := c.Locals("body").(*CreateUserRequest)
	req = *body

	response, err := h.service.CreateUser(req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		h.logger.Error("failed to create user",
			"error", err,
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// ChangePassword handles password change (self-service)
// @Summary Change user password
// @Description Change the password for the authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequest true "Password change request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/password/change [post]
func (h *Handler) ChangePassword(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	var req ChangePasswordRequest

	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("failed to parse request body",
			"error", err,
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := middleware.ValidateBody[ChangePasswordRequest](func(c *fiber.Ctx) error {
		return nil
	})(c); err != nil {
		return err
	}

	body := c.Locals("body").(*ChangePasswordRequest)
	req = *body

	if err := h.service.ChangePassword(userID, req); err != nil {
		if strings.Contains(err.Error(), "incorrect") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Current password is incorrect",
			})
		}
		h.logger.Error("failed to change password",
			"error", err,
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to change password",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password changed successfully",
	})
}

// GetUserInfo handles getting user information
// @Summary Get user information
// @Description Get information about the authenticated user
// @Tags auth
// @Produce json
// @Success 200 {object} UserInfoResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/me [get]
func (h *Handler) GetUserInfo(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	userInfo, err := h.service.GetUserInfo(userID)
	if err != nil {
		h.logger.Error("failed to get user info",
			"error", err,
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user information",
		})
	}

	return c.JSON(userInfo)
}
