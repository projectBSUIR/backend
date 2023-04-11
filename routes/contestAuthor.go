package routes

import (
	"fiber-apis/models"
	"fiber-apis/types"
	"github.com/gofiber/fiber/v2"
)

func GetContests(c *fiber.Ctx) error {
	userId, err := models.GetUserId(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	ownContests, err := models.FetchAllContestsForAuthor(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"contests": ownContests,
	})
}

func SetCoach(c *fiber.Ctx) error {
	var userInfo models.UserInfo
	if err := c.BodyParser(&userInfo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	err := models.UpdateStatus(userInfo.Id, types.Coach)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
