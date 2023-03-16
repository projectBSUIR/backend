package models

import (
	"errors"
	"fiber-apis/databases"
	_ "github.com/go-sql-driver/mysql"
)

type UserStatus int

const (
	Admin       UserStatus = 2
	Coach       UserStatus = 1
	Participant UserStatus = 0
)

type User struct {
	ID       int        `json:"id"`
	Login    string     `json:"login"`
	Password string     `json:"password"`
	Email    string     `json:"email"`
	Status   UserStatus `json:"status"`
}

func (model *User) LogIn() error {
	res, err := databases.DataBase.Query("SELECT count(*) FROM `user` WHERE `login` = ? AND `password` = ?", model.Login, model.Password)
	if err != nil {
		return err
	}
	var count int
	res.Next()
	err = res.Scan(&count)
	if err != nil {
		return err
	}

	if count == 1 {
		return nil
	}
	return errors.New("wrong login or password")
}

func (model *User) Register() error {
	res, err := databases.DataBase.Query("SELECT count(*) FROM `user` WHERE `login` = ?", model.Login)
	if err != nil {
		return err
	}
	var count int
	res.Next()
	err = res.Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("user already exists")
	} else {
		row, err := databases.DataBase.Exec("INSERT INTO `user` (`login`, `password`, `email`, `status`) VALUES (?, ?, ?, ?);",
			model.Login, model.Password, model.Email, model.Status)
		if err != nil {
			return err
		}
		id, err := row.LastInsertId()
		if err != nil {
			return err
		}
		model.ID = int(id)
	}
	return nil
}
