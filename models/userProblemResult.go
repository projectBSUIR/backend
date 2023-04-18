package models

import (
	"database/sql"
	"fiber-apis/databases"
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
	Result      int64        `json:"result"`
	LastAttempt sql.NullTime `json:"last_attempt,omitempty"`
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
	res, err := databases.DataBase.Query("SELECT `result`, `last_attempt` FROM `userProblemResult` WHERE `user_id`= ? AND `problem_id` IN (SELECT `problem_id` FROM `problem` WHERE `contest_id`= ?)", userId, contestId)
	if err != nil {
		return nil, err
	}
	for res.Next() {
		var r ProblemResultInfo
		err := res.Scan(&r.Result, &r.LastAttempt)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}

func (userProblemResult *UserProblemResult) Exists() (bool, error) {
	res, err := databases.DataBase.Query("SELECT count(*) FROM `userProblemResult` WHERE `user_id`= ? AND `problem_id`= ?", userProblemResult.UserId, userProblemResult.ProblemId)
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

func AddUserProblemResultIfNotExists(userId int64, problemId int64) error {
	userProblemResult := UserProblemResult{
		Id:            0,
		Result:        0,
		AttemptsCount: 0,
		LastAttempt:   "2006-01-02 15:04:05",
		UserId:        userId,
		ProblemId:     problemId,
	}

	exists, err := userProblemResult.Exists()
	if err != nil {
		return err
	}
	if !exists {
		err = userProblemResult.Create()
		if err != nil {
			return err
		}
	}
	return nil
}
