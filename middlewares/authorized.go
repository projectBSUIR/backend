package middlewares

import (
	"fiber-apis/models"
	"fiber-apis/types"
	"github.com/gofiber/fiber/v2"
)

func Participant(c *fiber.Ctx) error {
	userStatus, err := models.GetUserStatus(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if userStatus == types.UnAuthorized {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if userStatus < types.Participant {
		return c.SendStatus(fiber.StatusForbidden)
	}

	return c.Next()
}

func Coach(c *fiber.Ctx) error {
	userStatus, err := models.GetUserStatus(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if userStatus == types.UnAuthorized {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if userStatus < types.Coach {
		return c.SendStatus(fiber.StatusForbidden)
	}

	return c.Next()
}

func Admin(c *fiber.Ctx) error {
	userStatus, err := models.GetUserStatus(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if userStatus == types.UnAuthorized {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if userStatus != types.Admin {
		return c.SendStatus(fiber.StatusForbidden)
	}

	return c.Next()
}

func TestMachine(c *fiber.Ctx) error {
	userStatus, err := models.CheckTestMachine(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if userStatus == types.UnAuthorized {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if userStatus != types.TestMachine {
		return c.SendStatus(fiber.StatusForbidden)
	}

	return c.Next()
}
