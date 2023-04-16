package models

import (
	"fiber-apis/databases"
	"fiber-apis/token"
	"github.com/gofiber/fiber/v2"
)

type ContestResult struct {
	Id             int `json:"id"`
	SolvedProblems int `json:"solved_problems"`
	Penalty        int `json:"penalty"`
	UserId         int `json:"user_id"`
	ContestId      int `json:"contest_id"`
}

func GetUserContestResult(ContestId int, c *fiber.Ctx) (ContestResult, error) {
	var resModel ContestResult
	UserId, err := token.GetUserId(c)
	if err != nil {
		return resModel, err
	}
	res, err := databases.DataBase.Query("SELECT * FROM `contest_result` WHERE `user_id`= ? AND `contest_id`= ?", UserId, ContestId)
	if err != nil {
		_, err := databases.DataBase.Query("ROLLBACK")
		if err != nil {
			return resModel, err
		}
	}
	res.Next()
	err = res.Scan(&resModel.Id, &resModel.SolvedProblems, &resModel.Penalty, &resModel.UserId, &resModel.ContestId)
	if err != nil {
		return resModel, err
	}
	return resModel, nil
}
