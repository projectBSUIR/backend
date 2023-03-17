package routes

import "github.com/gofiber/fiber/v2"

func CheckHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON("Hello, authorized user")
}
