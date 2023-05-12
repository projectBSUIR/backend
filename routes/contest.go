package routes

import (
	"encoding/json"
	"fiber-apis/models"
	"fiber-apis/types"
	"fiber-apis/zipper"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

type ProblemData struct {
	ContestId int `json:"contestId"`
}

func (problem *ProblemData) GetProblemModel(testSet []byte, properties []byte, checker []byte) models.Problem {
	return models.Problem{
		Id:         0,
		ContestId:  problem.ContestId,
		TestSet:    testSet,
		Properties: properties,
		Checker:    checker,
	}
}

type ContestData struct {
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
	Duration  int64     `json:"duration"`
}

func (contest *ContestData) GetContestModel() models.Contest {
	return models.Contest{
		Id:        0,
		Name:      contest.Name,
		StartTime: contest.StartTime.UTC().Format("2006-01-02 15:04:05"),
		Duration:  contest.Duration,
	}
}

func AddProblem(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("Problem")

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	file, err := fileHeader.Open()
	defer file.Close()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	files, err := zipper.ExtractAllInOrder(
		file,
		[]string{"tests/", "check.cpp", "statements/russian/problem-properties.json"},
		[]string{"tests.zip", "checker.cpp", "problem-properties.json"},
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	var problemData ProblemData
	err = json.Unmarshal([]byte(c.FormValue("Contest")), &problemData)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	problem := problemData.GetProblemModel(
		files[0],
		files[2],
		files[1],
	)

	err = problem.Create()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func CreateContest(c *fiber.Ctx) error {
	var contest ContestData

	if err := c.BodyParser(&contest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	contestModel := contest.GetContestModel()
	contestInfo, err := contestModel.Create(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(contestInfo)
}

func ViewContests(c *fiber.Ctx) error {
	contests, err := models.FetchAllContests()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"contests": contests,
	})
}

func ViewProblems(c *fiber.Ctx) error {
	contestId, err := strconv.Atoi(c.Params("contestId"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	var contest models.Contest

	err = contest.FetchContest(contestId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	userStatus, err := models.GetUserStatus(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if !contest.NotStarted() || userStatus == types.Admin {
		problems, err := models.GetProblemsFromContest(contestId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"problems": problems,
		})
	}

	return c.SendStatus(fiber.StatusForbidden)
}

func ExtractProblemTests(c *fiber.Ctx) error {
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
	base64Testset, err := models.GetTestset(testingInfo.ProblemId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	binaryTestset, err := types.EncodeBytesToBinary(base64Testset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).SendString(string(binaryTestset))
}
