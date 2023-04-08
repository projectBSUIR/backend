package models

type ContestResult struct {
	Id             int `json:"id"`
	SolvedProblems int `json:"solved_problems"`
	Penalty        int `json:"penalty"`
	UserId         int `json:"user_id"`
	ContestId      int `json:"contest_id"`
}
