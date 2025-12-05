package auth

import (
	"innotech/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers authentication routes
func RegisterRoutes(app *fiber.App, handler *Handler, jwtMiddleware *middleware.KeycloakJWTMiddleware) {
	auth := app.Group("/api/auth")

	protected := auth.Group("", jwtMiddleware.RequireAuth())
	{
		protected.Get("/me", handler.GetUserInfo)
		protected.Post("/password/change", handler.ChangePassword)
	}

	admin := auth.Group("/users", jwtMiddleware.RequireAuth(), jwtMiddleware.RequireRole("admin"))
	{
		admin.Post("", handler.CreateUser)
	}
}
