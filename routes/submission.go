package routes

import (
	"bytes"
	"fiber-apis/models"
	"fiber-apis/token"
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"strconv"
)

const MAX_SOLUTION_SIZE int64 = 1024 * 1024 * 1024 * 2

func SubmitSolution(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("Solution")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	buf := bytes.NewBuffer(nil)
	written, err := io.Copy(buf, file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if written > MAX_SOLUTION_SIZE {
		return c.SendStatus(fiber.StatusRequestEntityTooLarge)
	}

	var submission models.Submission
	submission.Solution = buf.Bytes()
	submission.UserId, err = token.GetUserId(c)
	if err != nil {
		if err.Error() == "refresh_token is expired" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var problemId int64
	problemId, err = strconv.ParseInt(c.FormValue("ProblemId"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	submission.ProblemId = problemId
	submission.SetDefaultValues()
	err = submission.Create()
	log.Println(submission)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var testingQueue models.TestingQueue
	testingQueue.SubmissionId = submission.Id
	err = testingQueue.AddSubmissionToQueue()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
