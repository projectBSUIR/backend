package routes

import (
	"encoding/json"
	"fiber-apis/models"
	"fiber-apis/zipper"
	"github.com/gofiber/fiber/v2"
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
	//log.Println(contest.StartTime.Format("2006-01-02 15:04:05"))
	return models.Contest{
		Id:        0,
		Name:      contest.Name,
		StartTime: contest.StartTime.Format("2006-01-02 15:04:05"),
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
		[]string{"tests.zip", "checker.zip", "problem-properties.json"},
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
	err := contestModel.Create()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"contestId": contestModel.Id,
	})
}
