package models

import (
	"fiber-apis/databases"
	"github.com/gofiber/fiber/v2"
)

type ContestAuthor struct {
	Id        int64 `json:"id"`
	UserId    int64 `json:"user_id"`
	ContestId int64 `json:"contest_id"`
}

func CreateContestAuthor(Id int64, userId int64, contestId int64) ContestAuthor {
	return ContestAuthor{
		Id:        Id,
		UserId:    userId,
		ContestId: contestId,
	}
}

func SetAuthorOfContest(contestId int64, c *fiber.Ctx) error {
	id, err := GetUserId(c)
	if err != nil {
		return err
	}
	contestAuthor := ContestAuthor{
		Id:        0,
		UserId:    id,
		ContestId: contestId,
	}
	return contestAuthor.Create()
}

func (contestAuthor *ContestAuthor) IsAuthorOfContest() (bool, error) {
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

func (contestAuthor *ContestAuthor) Create() error {
	row, err := databases.DataBase.Exec("INSERT INTO `contestAuthor` (`user_id`, `contest_id`) VALUES (?, ?)", contestAuthor.UserId, contestAuthor.ContestId)
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
	contestAuthor.Id = id
	return nil
}
