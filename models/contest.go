package models

import (
	"fiber-apis/databases"
	_ "github.com/go-sql-driver/mysql"
)

type Contest struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	StartTime string `json:"start_time"`
	Duration  int64  `json:"duration"`
}

func (contest *Contest) Create() error {
	row, err := databases.DataBase.Exec("INSERT INTO `contest` (`contest_name`, `start_time`, `duration`) VALUES (?, ?, ?)",
		contest.Name, contest.StartTime, contest.Duration) // Надо проверить как оно здесь записывает дату и время
	if err != nil {
		return err
	}
	id, err := row.LastInsertId()
	if err != nil {
		return err
	}
	contest.Id = int(id)
	return nil
}
