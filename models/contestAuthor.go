package models

import "fiber-apis/databases"

type ContestAuthor struct {
	Id        int64 `json:"id"`
	UserId    int64 `json:"user_id"`
	ContestId int64 `json:"contest_id"`
}

func (contestAuthor *ContestAuthor) isAuthorOfContest() (bool, error) {
	rows, err := databases.DataBase.Query("SELECT count(*) FROM `contestAuthor` WHERE `user_id`=? AND `contest_id`=?", contestAuthor.UserId, contestAuthor.ContestId)
	if err != nil {
		return false, err
	}
	var count int
	rows.Next()
	rows.Scan(&count)

	if count == 1 {
		return true, nil
	}
	return false, nil
}
