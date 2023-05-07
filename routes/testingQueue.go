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

func ExtractFilesForTesting(c *fiber.Ctx) error {
	var testingInfo models.TestingIdsInfo
	if err := c.BodyParser(&testingInfo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	testingFilesInfo, err := models.GetFilesForTestingSubmission(testingInfo.SubmissionId, testingInfo.ProblemId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(testingFilesInfo)
}
