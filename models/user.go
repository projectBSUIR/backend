package models

import (
	"errors"
	"fiber-apis/databases"
	"fiber-apis/token"
	"fiber-apis/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"time"
)

type User struct {
	ID       int64            `json:"id"`
	Login    string           `json:"login"`
	Password string           `json:"password"`
	Email    string           `json:"email"`
	Status   types.UserStatus `json:"status"`
}

type UserInfo struct {
	Id int64 `json:"id"`
}

type TestMachineBody struct {
	Login    string      `json:"login"`
	Password string      `json:"password"`
	Payload  interface{} `json:"payload"`
}

func GetJWTClaim(model *User, expirationTime time.Time) token.JWTClaim {
	return token.JWTClaim{
		Id:     model.ID,
		User:   model.Login,
		Status: model.Status,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
}

func (model *User) SetStatus(s string) {
	model.Status = types.GetStatusByString(s)
}

func (model *User) LogIn() error {
	res, err := databases.DataBase.Query("SELECT id, email, status FROM `user` WHERE `login` = ? AND `password` = ?", model.Login, model.Password)
	if err != nil {
		return err
	}
	var count int = 0
	var status string
	for res.Next() {
		count++
		if count == 1 {
			res.Scan(&model.ID, &model.Email, &status)
		}
	}
	model.SetStatus(status)
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
			_, err := databases.DataBase.Query("ROLLBACK")
			if err != nil {
				return err
			}
			return err
		}
		id, err := row.LastInsertId()
		if err != nil {
			_, err := databases.DataBase.Query("ROLLBACK")
			if err != nil {
				return err
			}
			return err
		}
		model.ID = id
	}
	return nil
}

func CheckTestMachine(c *fiber.Ctx) (types.UserStatus, error) {
	var body TestMachineBody
	if err := c.BodyParser(&body); err != nil {
		return types.UnAuthorized, err
	}
	log.Println(body)
	user := User{
		Login:    body.Login,
		Password: body.Password,
	}
	err := user.LogIn()
	if err != nil {
		return types.UnAuthorized, err
	}
	return user.Status, nil
}

func GetTestMachineRequestPayload(c *fiber.Ctx) (interface{}, error) {
	var body TestMachineBody
	if err := c.BodyParser(&body); err != nil {
		return 0, err
	}
	return body.Payload, nil
}

func UpdateStatus(userId int64, status types.UserStatus) error {
	_, err := databases.DataBase.Exec("UPDATE `user` SET `status`=? WHERE `id`=?", status, userId)
	if err != nil {
		prevErr := err
		_, err := databases.DataBase.Query("ROLLBACK")
		if err != nil {
			return err
		}
		return prevErr
	}
	return nil
}

func GetUserStatus(c *fiber.Ctx) (types.UserStatus, error) {
	userInfo, err := token.GetClaimsFromTokens(c)
	if err != nil {
		return types.UnAuthorized, err
	}
	return userInfo.Status, nil
}

func GetUserId(c *fiber.Ctx) (int64, error) {
	userInfo, err := token.GetClaimsFromTokens(c)
	if err != nil {
		return 0, err
	}

	return userInfo.Id, nil
}

func GetLoginById(userId int64) (string, error) {
	var login string
	log, err := databases.DataBase.Query("SELECT `login` FROM `user` WHERE `id`= ?", userId)
	if err != nil {
		return "", err
	}
	defer log.Close()
	log.Next()
	err = log.Scan(&login)
	if err != nil {
		return "", err
	}
	return login, nil
}
