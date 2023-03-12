package routes

import "github.com/gofiber/fiber/v2"

func RegisterHandler(c *fiber.Ctx) error {

	return c.Status(200).JSON("Register successful")
}
