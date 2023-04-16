package models

import (
	"fiber-apis/databases"
	"github.com/gofiber/fiber/v2"
)

type UserProblemResult struct {
	Id            int    `json:"id"`
	Result        int    `json:"result"`
	AttemptsCount int    `json:"attempts_count"`
	LastAttempt   string `json:"last_attempt"`
	UserId        int    `json:"user_id"`
	ContestId     int    `json:"contest_id"`
}

func GetResultsFromContest(ContestId int, c *fiber.Ctx) ([]UserProblemResult, error) {
	var result []UserProblemResult
	UserId, err := GetUserId(c)
	if err != nil {
		return result, err
	}
	res, err := databases.DataBase.Query("SELECT * FROM `userProblemResult` WHERE `user_id` = ? AND `problem_id` IN (SELECT `problem_id` FROM `problem` WHERE `contest_id` = ?)", UserId, ContestId)
	if err != nil {
		_, err := databases.DataBase.Query("ROLLBACK")
		if err != nil {
			return result, err
		}
	}
	for res.Next() {
		var tmp UserProblemResult
		err := res.Scan(&tmp.Id, &tmp.Result, &tmp.AttemptsCount, &tmp.LastAttempt, &tmp.UserId, &tmp.ContestId)
		if err != nil {
			return result, err
		}
		result = append(result, tmp)
	}
	return result, nil
}
