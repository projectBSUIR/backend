package routes

import (
	"bytes"
	"fiber-apis/models"
	"github.com/gofiber/fiber/v2"
	"io"
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
	submission.UserId, err = models.GetUserId(c)
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
	contestId, err := strconv.ParseInt(c.FormValue("ContestId"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	submission.ProblemId = problemId
	submission.SetDefaultValues()
	err = submission.Create()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = models.AddContestResultIfNotExists(submission.UserId, contestId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	err = models.AddUserProblemResultIfNotExists(submission.UserId, submission.ProblemId)
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

func GetSubmissions(c *fiber.Ctx) error {
	problemId, err := strconv.ParseInt(c.Params("problemId"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	userId, err := models.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	submissionsInfo, err := models.GetSubmissionsByProblem(userId, problemId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"submissions": submissionsInfo,
	})
}

func GetAllSubmissions(c *fiber.Ctx) error {
	userId, err := models.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	contestId, err := strconv.ParseInt(c.Params("contestId"), 10, 64)
	contestAuthor := models.ContestAuthor{
		Id:        0,
		UserId:    userId,
		ContestId: contestId,
	}
	isAuthor, err := contestAuthor.IsAuthorOfContest()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if !isAuthor {
		return c.SendStatus(fiber.StatusForbidden)
	}
	submissionsInfo, err := models.GetSubmissions()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"submissions": submissionsInfo,
	})
}
