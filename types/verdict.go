package types

type VerdictInfo struct {
	Status string `json:"status"`
	Time   int64  `json:"time"`
	Memory int64  `json"memory"`
}

type TestingVerdict struct {
	SubmissionId int64       `json:"submission_id"`
	ProblemId    int64       `json:"problem_id"`
	UserId       int64       `json:"user_id"`
	Verdict      VerdictInfo `json:"verdict"`
}
