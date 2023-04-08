package models

import (
	"fiber-apis/databases"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Contest struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	StartTime string `json:"start_time"`
	Duration  int64  `json:"duration"`
}

func (contest *Contest) NotStarted() bool {
	timeNow := time.Now()
	startTimeContest, _ := time.Parse(time.RFC3339, contest.StartTime)
	return timeNow.Before(startTimeContest)
}

func (contest *Contest) Create() error {
	row, err := databases.DataBase.Exec("INSERT INTO `contest` (`contest_name`, `start_time`, `duration`) VALUES (?, ?, ?)",
		contest.Name, contest.StartTime, contest.Duration) // Надо проверить как оно здесь записывает дату и время
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
	contest.Id = int(id)
	return nil
}

func GetAllContests() ([]Contest, error) {
	rows, err := databases.DataBase.Query("SELECT * FROM `contest`")
	if err != nil {
		_, err := databases.DataBase.Query("ROLLBACK")
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	contests := make([]Contest, 0)
	for rows.Next() {
		var contest Contest
		rows.Scan(&contest.Id, &contest.Name, &contest.StartTime, &contest.Duration)
		contests = append(contests, contest)
	}
	return contests, nil
}

func (contest *Contest) GetContest(contestId int) error {
	row, err := databases.DataBase.Query("SELECT `contest_name`, `start_time`, `duration` FROM `contest` WHERE id = ?", contestId)
	if err != nil {
		return err
	}

	row.Next()
	err = row.Scan(&contest.Name, &contest.StartTime, &contest.Duration)
	if err != nil {
		return err
	}

	return nil
}
