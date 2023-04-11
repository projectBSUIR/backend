package models

type UserProblemResult struct {
	Id                   int    `json:"id"`
	Result               int    `json:"result"`
	AttemptsCount        int    `json:"attempts_count"`
	LastAttempt          string `json:"last_attempt"`
	UserProblemResultCol string `json:"user_problem_result_col"`
	UserId               int    `json:"user_id"`
	ContestId            int    `json:"contest_id"`
}
