package middlewares

import (
	"fiber-apis/models"
	"fiber-apis/token"
	"github.com/gofiber/fiber/v2"
)

func Participant(c *fiber.Ctx) error {
	userStatus, err := token.GetUserStatus(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if userStatus == models.UnAuthorized {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.Next()
}

func Coach(c *fiber.Ctx) error {
	userStatus, err := token.GetUserStatus(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if userStatus == models.UnAuthorized {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if userStatus < models.Admin {
		return c.SendStatus(fiber.StatusForbidden)
	}

	return c.Next()
}

func Admin(c *fiber.Ctx) error {
	userStatus, err := token.GetUserStatus(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if userStatus == models.UnAuthorized {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if userStatus != models.Admin {
		return c.SendStatus(fiber.StatusForbidden)
	}

	return c.Next()
}
