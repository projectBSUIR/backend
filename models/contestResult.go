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

func CreateContestResult(Id int64, Penalty int64, SolvedTasks int64, userId int64, contestId int64) ContestResult {
	return ContestResult{
		Id:          Id,
		Penalty:     Penalty,
		SolvedTasks: SolvedTasks,
		UserId:      userId,
		ContestId:   contestId,
	}
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

func (contestResult *ContestResult) GetContestResultInfo() error {
	row, err := databases.DataBase.Query("SELECT `id`, `penalty`, `solvedTasks` FROM `contestResult` WHERE `user_id`=? AND `contest_id`=?", contestResult.UserId, contestResult.ContestId)
	if err != nil {
		return err
	}
	row.Next()
	err = row.Scan(&contestResult.Id, &contestResult.Penalty, &contestResult.SolvedTasks)
	return err
}

func (contestResult *ContestResult) UpdateContestResultInfo() error {
	_, err := databases.DataBase.Exec("UPDATE `contestResult` SET `penalty`=?, `solvedTasks`=? WHERE `id`=?", contestResult.Penalty, contestResult.SolvedTasks, contestResult.Id)
	if err != nil {
		_, nerr := databases.DataBase.Exec("ROLLBACK")
		if nerr != nil {
			return nerr
		}
		return err
	}
	return nil
}

func AddContestResultIfNotExists(userId int64, contestId int64) error {
	contestResult := CreateContestResult(0, 0, 0, userId, contestId)
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

func GetContestIdByProblemId(problemId int64) (int64, error) {
	row, err := databases.DataBase.Query("SELECT `contest_id` FROM `problem` WHERE `id`=? ", problemId)
	if err != nil {
		return 0, err
	}
	var contestId int64
	row.Next()
	err = row.Scan(&contestId)
	if err != nil {
		return 0, err
	}
	return contestId, nil
}

func UpdateUserContestResult(userId int64, problemId int64, attemptsCount int64) error {
	contestId, err := GetContestIdByProblemId(problemId)
	if err != nil {
		return err
	}
	contestResult := CreateContestResult(0, 0, 0, userId, contestId)
	err = contestResult.GetContestResultInfo()
	if err != nil {
		return err
	}
	contestStartTime, err := GetStartTimeOfContest(contestId)
	if err != nil {
		return err
	}
	submitTime, err := GetSubmitTime(userId, problemId, attemptsCount+1)
	if err != nil {
		return err
	}

	contestResult.Penalty += attemptsCount*20 + int64(submitTime.Sub(contestStartTime).Minutes())
	contestResult.SolvedTasks++

	return contestResult.UpdateContestResultInfo()
}
