package routes

import (
	"fiber-apis/models"
	"github.com/gofiber/fiber/v2"
)

type UserResult struct {
	Login          string                     `json:"login"`
	ContestResult  int64                      `json:"contest_result"`
	Penalty        int64                      `json:"penalty"`
	ProblemsResult []models.ProblemResultInfo `json:"problems_result"`
}

func CreateTable(contestId int64) ([]UserResult, error) {
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
			ContestResult:  contestResult.SolvedTasks,
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
	table, err := CreateTable(contestId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"results": table,
	})
}
