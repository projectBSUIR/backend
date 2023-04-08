package models

import (
	"fiber-apis/databases"
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

func (submission *Submission) Create() error {
	row, err := databases.DataBase.Exec("INSERT INTO `submission` "+
		"(`id`, `solution`, `submit_time`, `verdict`, `problem_id`, `user_id`) "+
		"VALUES (?, ?, ?, ?, ?, ?)",
		submission.Id, submission.Solution, submission.SubmitTime,
		submission.Verdict, submission.ProblemId, submission.UserId)

	if err != nil {
		_, err := databases.DataBase.Query("ROLLBACK")
		if err != nil {
			return err
		}
		return err
	}
	id, err := row.LastInsertId()
	if err != nil {
		_, err := databases.DataBase.Query("ROLLBACK")
		if err != nil {
			return err
		}
	}
	submission.Id = id
	return nil
}

func GetSubmissionsByUser(userId int64) ([]SubmissionInfo, error) {
	rows, err := databases.DataBase.Query("SELECT `id`, `submit_time`, `verdict`, `problem_id` FROM `submission` WHERE `user_id` = ?", userId)
	if err != nil {
		return nil, err
	}
	submissions := make([]SubmissionInfo, 0)
	for rows.Next() {
		var submission SubmissionInfo
		err := rows.Scan(&submission.Id, &submission.SubmitTime, &submission.Verdict, &submission.ProblemId)
		submission.UserId = userId
		submissions = append(submissions, submission)
		if err != nil {
			return nil, err
		}
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
		submissions = append(submissions, submission)
		if err != nil {
			return nil, err
		}
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

func UpdateSubmissionVerdict(submissionId int64, newVerdict string) error {
	_, err := databases.DataBase.Exec("UPDATE `submission` SET `verdict`=? WHERE `id`=?", newVerdict, submissionId)
	if err != nil {
		_, err := databases.DataBase.Query("ROLLBACK")
		if err != nil {
			return err
		}
		return err
	}
	return nil
}
