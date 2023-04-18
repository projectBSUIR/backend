package models

import (
	"encoding/json"
	"fiber-apis/databases"
	"github.com/gofiber/fiber/v2"
	"time"
)

type SubmissionInfo struct {
	Id         int64  `json:"id"`
	SubmitTime string `json:"submit_time"`
	Verdict    string `json:"verdict"`
	ProblemId  int64  `json:"problem_id"`
	UserId     int64  `json:"user_id"`
}

type Submission struct {
	Id         int64  `json:"id"`
	Solution   []byte `json:"solution"`
	SubmitTime string `json:"submit_time"`
	Verdict    string `json:"verdict"`
	ProblemId  int64  `json:"problem_id"`
	UserId     int64  `json:"user_id"`
}

func CreateVerdict(status string, log string, time int64, memory int64) string {
	verdict, _ := json.Marshal(fiber.Map{
		"Status": status,
		"Log":    log,
		"Time":   time,
		"Memory": memory,
	})
	return string(verdict)
}

func (submission *Submission) SetDefaultValues() {
	submission.Verdict = CreateVerdict("Pending", "", 0, 0)
	submission.SubmitTime = time.Now().Format(time.RFC3339)
}

func (submission *Submission) Create() error {
	row, err := databases.DataBase.Exec("INSERT INTO `submission` (`id`, `solution`, `submit_time`, `verdict`, `problem_id`, `user_id`) VALUES (?, ?, ?, ?, ?, ?)",
		submission.Id, submission.Solution, submission.SubmitTime,
		submission.Verdict, submission.ProblemId, submission.UserId)

	if err != nil {
		var prevErr error = err
		_, err := databases.DataBase.Query("ROLLBACK")
		if err != nil {
			return err
		}
		return prevErr
	}
	id, err := row.LastInsertId()
	if err != nil {
		var prevErr error = err
		_, err := databases.DataBase.Query("ROLLBACK")
		if err != nil {
			return err
		}
		return prevErr
	}
	submission.Id = id
	return nil
}

func GetSubmissionsByProblem(userId int64, problemId int64) ([]SubmissionInfo, error) {
	rows, err := databases.DataBase.Query("SELECT `id`, `submit_time`, `verdict` FROM `submission` WHERE `user_id` = ? AND`problem_id` = ?", userId, problemId)
	if err != nil {
		return nil, err
	}
	submissions := make([]SubmissionInfo, 0)
	for rows.Next() {
		var submission SubmissionInfo
		err := rows.Scan(&submission.Id, &submission.SubmitTime, &submission.Verdict)
		if err != nil {
			return nil, err
		}

		submission.ProblemId = problemId
		submission.UserId = userId

		submissions = append(submissions, submission)
	}
	return submissions, nil
}

func GetSubmissions() ([]SubmissionInfo, error) {
	rows, err := databases.DataBase.Query("SELECT `id`, `submit_time`, `verdict`, `problem_id`, `user_id` FROM `submission`")
	if err != nil {
		return nil, err
	}
	submissions := make([]SubmissionInfo, 0)
	for rows.Next() {
		var submission SubmissionInfo
		err := rows.Scan(&submission.Id, &submission.SubmitTime, &submission.Verdict, &submission.ProblemId, &submission.UserId)
		if err != nil {
			return nil, err
		}

		submissions = append(submissions, submission)
	}
	return submissions, nil
}

func GetSolutionBySubmissionId(submissionId int64) ([]byte, error) {
	rows, err := databases.DataBase.Query("SELECT `solution` FROM `submission` WHERE `id`=?", submissionId)
	if err != nil {
		return nil, err
	}
	var solution []byte
	for rows.Next() {
		err := rows.Scan(&solution)
		if err != nil {
			return nil, err
		}
	}
	return solution, nil
}

func UpdateSubmissionVerdict(submissionId int64, newVerdict fiber.Map) error {
	_, err := databases.DataBase.Exec("UPDATE `submission` SET `verdict`=? WHERE `id`=?", newVerdict, submissionId)
	if err != nil {
		var prevErr error = err
		_, err := databases.DataBase.Query("ROLLBACK")
		if err != nil {
			return err
		}
		return prevErr
	}
	return nil
}
