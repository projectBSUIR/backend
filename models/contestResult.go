package models

import (
	"fiber-apis/databases"
)

type ContestResult struct {
	Id             int `json:"id"`
	SolvedProblems int `json:"solved_problems"`
	Penalty        int `json:"penalty"`
	UserId         int `json:"user_id"`
	ContestId      int `json:"contest_id"`
}

func GetUserContestResult(userId int64, contestId int64) (ContestResult, error) {
	var resModel ContestResult
	res, err := databases.DataBase.Query("SELECT * FROM `contestResult` WHERE `user_id`= ? AND `contest_id`= ?", userId, contestId)
	if err != nil {
		_, nerr := databases.DataBase.Query("ROLLBACK")
		if nerr != nil {
			return resModel, nerr
		}
		return resModel, err
	}
	res.Next()
	err = res.Scan(&resModel.Id, &resModel.SolvedProblems, &resModel.Penalty, &resModel.UserId, &resModel.ContestId)
	if err != nil {
		return resModel, err
	}
	return resModel, nil
}
