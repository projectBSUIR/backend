package models

import (
	"fiber-apis/databases"
)

type ContestResult struct {
	Id          int64 `json:"id"`
	SolvedTasks int64 `json:"solvedTasks"`
	Penalty     int64 `json:"penalty"`
	UserId      int64 `json:"user_id"`
	ContestId   int64 `json:"contest_id"`
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
	err = res.Scan(&resModel.Id, &resModel.SolvedTasks, &resModel.Penalty, &resModel.UserId, &resModel.ContestId)
	if err != nil {
		return resModel, err
	}
	return resModel, nil
}

func (contestResult *ContestResult) Exists() (bool, error) {
	res, err := databases.DataBase.Query("SELECT count(*) FROM `contestResult` WHERE `user_id`= ? AND `contest_id`= ?", contestResult.UserId, contestResult.ContestId)
	if err != nil {
		return false, err
	}
	var count int64
	res.Next()
	err = res.Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

func (contestResult *ContestResult) Create() error {
	_, err := databases.DataBase.Exec("INSERT INTO `contestResult` (`solvedTasks`,`penalty`,`user_id`,`contest_id`) VALUES (?, ?, ?, ?)",
		contestResult.SolvedTasks, contestResult.Penalty,
		contestResult.UserId, contestResult.ContestId)

	if err != nil {
		_, nerr := databases.DataBase.Query("ROLLBACK")
		if nerr != nil {
			return nerr
		}
		return err
	}
	return nil
}

func AddContestResultIfNotExists(userId int64, contestId int64) error {
	contestResult := ContestResult{
		Id:          0,
		Penalty:     0,
		SolvedTasks: 0,
		UserId:      userId,
		ContestId:   contestId,
	}
	exists, err := contestResult.Exists()
	if err != nil {
		return err
	}
	if !exists {
		err = contestResult.Create()
		if err != nil {
			return err
		}
	}
	return nil
}
