package models

import (
	"database/sql"
	"errors"
	"fiber-apis/databases"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type UserProblemResult struct {
	Id            int64  `json:"id"`
	Result        int64  `json:"result"`
	AttemptsCount int64  `json:"attempts_count"`
	LastAttempt   string `json:"last_attempt"`
	UserId        int64  `json:"user_id"`
	ProblemId     int64  `json:"problem_id"`
}

type ProblemResultInfo struct {
	Result        int64        `json:"result"`
	AttemptsCount int64        `json:"attempts_count"`
	LastAttempt   sql.NullTime `json:"last_attempt,omitempty"`
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
		err := res.Scan(&tmp.Id, &tmp.Result, &tmp.AttemptsCount, &tmp.LastAttempt, &tmp.UserId, &tmp.ProblemId)
		if err != nil {
			return result, err
		}
		result = append(result, tmp)
	}
	return result, nil
}

func GetProblemsStatus(userId int64, contestId int64) ([]ProblemResultInfo, error) {
	var result []ProblemResultInfo
	res, err := databases.DataBase.Query("SELECT `result`, `attempts_count`, `last_attempt` FROM `userProblemResult` WHERE `user_id`= ? AND `problem_id` IN (SELECT `problem_id` FROM `problem` WHERE `contest_id`= ?)", userId, contestId)
	if err != nil {
		return nil, err
	}
	for res.Next() {
		var r ProblemResultInfo
		err := res.Scan(&r.Result, &r.AttemptsCount, &r.LastAttempt)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}

func (userProblemResult *UserProblemResult) Exists() (bool, error) {
	res, err := databases.DataBase.Query("SELECT `id`, `result`, `attempts_count` FROM `userProblemResult` WHERE `user_id`= ? AND `problem_id`= ?", userProblemResult.UserId, userProblemResult.ProblemId)
	if err != nil {
		return false, err
	}
	var count int64
	for res.Next() {
		count++
		err = res.Scan(&userProblemResult.Id, &userProblemResult.Result, &userProblemResult.AttemptsCount)
		if err != nil {
			return false, err
		}
	}
	return count == 1, nil
}

func (userProblemResult *UserProblemResult) Create() error {
	row, err := databases.DataBase.Exec("INSERT INTO `userProblemResult` (`result`, `attempts_count`, `last_attempt`, `user_id`, `problem_id`) VALUES (?, ?, ?, ?, ?)",
		userProblemResult.Result, userProblemResult.AttemptsCount,
		userProblemResult.LastAttempt,
		userProblemResult.UserId, userProblemResult.ProblemId)

	if err != nil {
		_, nerr := databases.DataBase.Query("ROLLBACK")
		if nerr != nil {
			return nerr
		}
		return err
	}
	userProblemResult.Id, err = row.LastInsertId()
	if err != nil {
		_, nerr := databases.DataBase.Query("ROLLBACK")
		if nerr != nil {
			return nerr
		}
		return err
	}
	return nil
}

func (userProblemResult *UserProblemResult) UpdateLastAttempt() error {
	if userProblemResult.Result == 1 {
		return nil
	}
	_, err := databases.DataBase.Exec("UPDATE `userProblemResult` SET `last_attempt`=? WHERE `user_id`=? AND `problem_id`=?",
		userProblemResult.LastAttempt,
		userProblemResult.UserId, userProblemResult.ProblemId)
	if err != nil {
		_, nerr := databases.DataBase.Query("ROLLBACK")
		if nerr != nil {
			return nerr
		}
		return err
	}
	return nil
}

func GetAttemptsCount(userId int64, problemId int64) (int64, error) {
	rows, err := databases.DataBase.Query("SELECT `attempts_count` FROM `userProblemResult` WHERE `user_id`=? AND `problem_id`=?",
		userId, problemId)
	if err != nil {
		return 0, err
	}
	var attemptsCount int64
	rows.Next()
	rows.Scan(&attemptsCount)
	return attemptsCount, nil
}

func AddUserProblemResultIfNotExists(userId int64, problemId int64, last_attempt string) error {
	userProblemResult := UserProblemResult{
		Id:            0,
		Result:        0,
		AttemptsCount: 0,
		LastAttempt:   last_attempt,
		UserId:        userId,
		ProblemId:     problemId,
	}

	exists, err := userProblemResult.Exists()
	if err != nil {
		return err
	}
	if !exists {
		err = userProblemResult.Create()
		return err
	}
	return nil
}

func GetUserProblemResult(userId int64, problemId int64) (int8, error) {
	rows, err := databases.DataBase.Query("SELECT `result` FROM `userProblemResult` WHERE `user_id`= ? AND `problem_id`= ?", userId, problemId)
	if err != nil {
		return 0, err
	}
	count := 0
	var result int8 = 0
	for rows.Next() {
		count++
		rows.Scan(&result)
	}
	if count != 1 {
		return 0, errors.New(fmt.Sprint("There is not one row with userId=", userId, " and problemId=", problemId))
	}
	return result, nil
}

func UpdateUserProblemResult(userId int64, problemId int64, new_result int8) error {
	result, err := GetUserProblemResult(userId, problemId)
	if err != nil {
		return err
	}
	if result == 1 {
		return nil
	}

	attemptsCount, err := GetAttemptsCount(userId, problemId)
	if err != nil {
		return err
	}

	if new_result == 1 {
		err = UpdateUserContestResult(userId, problemId, attemptsCount)
		if err != nil {
			return err
		}
	}

	attemptsCount++
	_, err = databases.DataBase.Exec("UPDATE `userProblemResult` SET `attempts_count`= ?, `result`= ? WHERE `user_id`= ? AND `problem_id`= ?", attemptsCount, new_result, userId, problemId)
	if err != nil {
		_, nerr := databases.DataBase.Exec("ROLLBACK")
		if nerr != nil {
			return nerr
		}
		return err
	}
	return nil
}
