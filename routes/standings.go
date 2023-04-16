package routes

import (
	"fiber-apis/models"
	"github.com/gofiber/fiber/v2"
)

type UserResult struct {
	Login          string `json:"login"`
	ContestResult  int    `json:"contest_result"`
	Penalty        int    `json:"penalty"`
	ProblemsResult []int  `json:"problems_result"`
}

func CreateTable(contestId int64, c *fiber.Ctx) ([]UserResult, error) {
	var table []UserResult
	contestParticipants, err := models.GetParticipantsIds(contestId)
	if err != nil {
		return nil, err
	}
	for _, userId := range contestParticipants {
		login, err := models.GetLoginById(userId)
		if err != nil {
			return nil, err
		}
		contestResult, err := models.GetUserContestResult(userId, contestId)
		if err != nil {
			return nil, err
		}
		userProblemResults, err := models.GetProblemsStatus(userId, contestId)
		if err != nil {
			return nil, err
		}
		table = append(table, UserResult{
			Login:          login,
			ContestResult:  contestResult.SolvedProblems,
			Penalty:        contestResult.Penalty,
			ProblemsResult: userProblemResults,
		})
	}
	return table, nil
}

func GetResultsTable(c *fiber.Ctx) error {
	var contestId int64
	err := c.BodyParser(&contestId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	table, err := CreateTable(contestId, c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"results": table,
	})
}
