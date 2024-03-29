package models

import (
	"encoding/json"
	"fiber-apis/databases"
	"github.com/gofiber/fiber/v2"
)

type Problem struct {
	Id         int    `json:"id"`
	ContestId  int    `json:"contestId"`
	TestSet    []byte `json:"testSet"`
	Properties []byte `json:"properties"`
	Checker    []byte `json:"checker"`
}

type ProblemInfo struct {
	Id         int       `json:"id"`
	Properties fiber.Map `json:"properties"`
}

func (problemInfo *ProblemInfo) SetInfo(problem Problem) error {
	err := json.Unmarshal(problem.Properties, &problemInfo.Properties)
	if err != nil {
		return err
	}

	problemInfo.Id = problem.Id
	return nil
}

func (problem *Problem) Create() error {
	id, err := databases.DataBase.Exec("INSERT INTO `problem` (`contest_id`, `testset`, `checker`, `problem_properties`) VALUES (?, ?, ?, ?)",
		problem.ContestId, problem.TestSet, problem.Checker, problem.Properties)
	if err != nil {
		return err
	}
	ID, err := id.LastInsertId()
	if err != nil {
		return err
	}
	problem.Id = int(ID)
	return nil
}

func GetProblemsFromContest(contestId int) ([]ProblemInfo, error) {
	rows, err := databases.DataBase.Query("SELECT `id`, `problem_properties` FROM `problem` WHERE `contest_id` = ?", contestId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	problems := make([]ProblemInfo, 0)
	for rows.Next() {
		var problem Problem
		rows.Scan(&problem.Id, &problem.Properties)

		var problemInfo ProblemInfo
		err := problemInfo.SetInfo(problem)
		if err != nil {
			return nil, err
		}

		problems = append(problems, problemInfo)
	}
	return problems, nil
}

func GetTestset(problemId int64) ([]byte, error) {
	res, err := databases.DataBase.Query("SELECT `testset` FROM `problem` WHERE `id`=?", problemId)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var testset []byte
	res.Next()
	err = res.Scan(&testset)
	return testset, err
}
