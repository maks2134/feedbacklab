package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ValidateBody[T any](next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body T

		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid JSON: " + err.Error(),
			})
		}

		if err := validate.Struct(&body); err != nil {
			errors := make(map[string]string)
			for _, e := range err.(validator.ValidationErrors) {
				errors[e.Field()] = e.Tag()
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"validation_errors": errors,
			})
		}

		c.Locals("body", &body)
		return next(c)
	}
}
