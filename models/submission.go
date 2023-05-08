package models

import (
	"encoding/json"
	"fiber-apis/databases"
	"fiber-apis/types"
	"github.com/gofiber/fiber/v2"
	"time"
)

type SubmissionInfo struct {
	Id         int64     `json:"id"`
	SubmitTime string    `json:"submit_time"`
	Verdict    fiber.Map `json:"verdict"`
	ProblemId  int64     `json:"problem_id"`
	UserId     int64     `json:"user_id"`
}

type Submission struct {
	Id         int64             `json:"id"`
	Solution   []byte            `json:"solution"`
	SubmitTime string            `json:"submit_time"`
	Verdict    types.VerdictInfo `json:"verdict"`
	ProblemId  int64             `json:"problem_id"`
	UserId     int64             `json:"user_id"`
}

type TestingIdsInfo struct {
	SubmissionId int64 `json:"submission_id"`
	ProblemId    int64 `json:"problem_id"`
}

type TestingFilesInfo struct {
	Solution          []byte    `json:"solution"`
	Testset           []byte    `json:"testset"`
	Checker           []byte    `json:"checker"`
	ProblemProperties fiber.Map `json:"problem_properties"`
}

func CreateVerdict(status string, time int64, memory int64) types.VerdictInfo {
	return types.VerdictInfo{
		Status: status,
		Time:   time,
		Memory: memory,
	}
}

func ConvertMapToString(verdict any) (string, error) {
	ret, err := json.Marshal(verdict)
	return string(ret), err
}

func ConvertToMap(sverdict string) fiber.Map {
	var verdict fiber.Map
	_ = json.Unmarshal([]byte(sverdict), &verdict)
	return verdict
}

func (submission *Submission) SetDefaultValues() {
	submission.Verdict = CreateVerdict("Pending", 0, 0)
	submission.SubmitTime = time.Now().UTC().Format("2006-01-02 15:04:05")
}

func (submission *Submission) Create() error {
	sverdict, err := ConvertMapToString(submission.Verdict)
	if err != nil {
		return err
	}
	row, err := databases.DataBase.Exec("INSERT INTO `submission` (`solution`, `submit_time`, `verdict`, `problem_id`, `user_id`) VALUES (?, ?, ?, ?, ?)",
		submission.Solution, submission.SubmitTime,
		sverdict, submission.ProblemId, submission.UserId)

	if err != nil {
		var prevErr error = err
		_, err := databases.DataBase.Exec("ROLLBACK")
		if err != nil {
			return err
		}
		return prevErr
	}
	id, err := row.LastInsertId()
	if err != nil {
		var prevErr error = err
		_, err := databases.DataBase.Exec("ROLLBACK")
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
		var stringVerdict string
		err := rows.Scan(&submission.Id, &submission.SubmitTime, &stringVerdict)
		if err != nil {
			return nil, err
		}
		submission.Verdict = ConvertToMap(stringVerdict)

		submission.ProblemId = problemId
		submission.UserId = userId

		submissions = append(submissions, submission)
	}
	return submissions, nil
}

func GetSubmissionsByContestId(contestId int64) ([]SubmissionInfo, error) {
	rows, err := databases.DataBase.Query("SELECT `id`, `submit_time`, `verdict`, `problem_id`, `user_id` FROM `submission` WHERE `problem_id` IN (SELECT `id` FROM `problem` WHERE `contest_id`= ?)", contestId)
	if err != nil {
		return nil, err
	}
	submissions := make([]SubmissionInfo, 0)
	for rows.Next() {
		var submission SubmissionInfo
		var sverdict string
		err := rows.Scan(&submission.Id, &submission.SubmitTime, &sverdict, &submission.ProblemId, &submission.UserId)
		submission.Verdict = ConvertToMap(sverdict)
		if err != nil {
			return nil, err
		}

		submissions = append(submissions, submission)
	}
	return submissions, nil
}

func UpdateSubmissionVerdict(submissionId int64, newVerdict fiber.Map) error {
	sverdict, err := ConvertMapToString(newVerdict)
	if err != nil {
		return err
	}
	_, err = databases.DataBase.Exec("UPDATE `submission` SET `verdict`=? WHERE `id`=?", sverdict, submissionId)
	if err != nil {
		var prevErr error = err
		_, err := databases.DataBase.Exec("ROLLBACK")
		if err != nil {
			return err
		}
		return prevErr
	}
	return nil
}

func GetFirstSubmissionFromTestingQueue() (TestingIdsInfo, error) {
	row, err := databases.DataBase.Query("SELECT `id`, `problem_id` FROM `submission` WHERE `id`= (SELECT `submission_id` FROM `testingQueue` ORDER BY `id` LIMIT 1)")
	if err != nil {
		return TestingIdsInfo{}, err
	}
	var testingIdsInfo TestingIdsInfo
	row.Next()
	err = row.Scan(&testingIdsInfo.SubmissionId, &testingIdsInfo.ProblemId)
	return testingIdsInfo, err
}

func GetFilesForTestingSubmission(submissionId int64, problemId int64) (TestingFilesInfo, error) {
	row, err := databases.DataBase.Query("SELECT `solution` FROM `submission` WHERE `id`= ?", submissionId)
	if err != nil {
		return TestingFilesInfo{}, err
	}
	var solutionInfo TestingFilesInfo
	row.Next()
	err = row.Scan(&solutionInfo.Solution)
	if err != nil {
		return TestingFilesInfo{}, err
	}

	row, err = databases.DataBase.Query("SELECT `testset`, `checker`, `problem_properties` FROM `problem` WHERE `id`= ?", problemId)
	if err != nil {
		return TestingFilesInfo{}, err
	}

	row.Next()
	var sproperties string

	err = row.Scan(&solutionInfo.Testset, &solutionInfo.Checker, &sproperties)
	if err != nil {
		return TestingFilesInfo{}, err
	}

	solutionInfo.ProblemProperties = ConvertToMap(sproperties)
	return solutionInfo, nil
}

func GetSubmitTime(userId int64, problemId int64, attemptId int64) (time.Time, error) {
	res, err := databases.DataBase.Query("SELECT `submit_time` FROM `submission` WHERE `user_id`=? AND `problem_id`=? ORDER BY `id` LIMIT ?", userId, problemId, attemptId)
	if err != nil {
		return time.Now(), err
	}
	var startTime time.Time
	var sstartTime string
	for i := int64(0); i < attemptId; i++ {
		res.Next()
	}
	err = res.Scan(&sstartTime)
	if err != nil {
		return time.Now(), err
	}
	startTime, err = time.Parse(time.RFC3339, sstartTime)
	return startTime, err
}

func GetUserIdFromSubmission(submissionId int64) (int64, error) {
	res, err := databases.DataBase.Query("SELECT `user_id` FROM `submission` WHERE `id`=?", submissionId)
	if err != nil {
		return 0, err
	}
	var userId int64
	res.Next()
	err = res.Scan(&userId)
	return userId, err
}
