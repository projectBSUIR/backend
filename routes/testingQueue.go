package routes

import (
	"fiber-apis/models"
	"github.com/gofiber/fiber/v2"
)

func ExtractSubmissionFromTestingQueue(c *fiber.Ctx) error {
	testingInfo, err := models.GetFirstSubmissionFromTestingQueue()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	err = models.DeleteSubmissionFromTestingQueue(testingInfo.SubmissionId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(testingInfo)
}
