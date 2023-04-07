package models

import (
	"fiber-apis/databases"
)

type Problem struct {
	Id         int    `json:"id"`
	ContestId  int    `json:"contestId"`
	TestSet    []byte `json:"testSet"`    //
	Properties []byte `json:"properties"` // Пока что не знаю правильно ли это
	Checker    []byte `json:"checker"`    //
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
