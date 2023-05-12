package routes

import (
	"encoding/hex"
	"encoding/json"
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
	payload, err := models.GetTestMachineRequestPayload(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	marshaledPayload, err := json.Marshal(payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err,
		})
	}
	var testingInfo models.TestingIdsInfo
	err = json.Unmarshal(marshaledPayload, &testingInfo)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err,
		})
	}
	testingFilesInfo, err := models.GetFilesForTestingSubmission(testingInfo.SubmissionId, testingInfo.ProblemId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"solution":           testingFilesInfo.Solution,
		"checker":            hex.EncodeToString(testingFilesInfo.Checker),
		"problem-properties": testingFilesInfo.ProblemProperties,
	})
}
