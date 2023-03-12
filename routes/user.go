package routes

import (
	"fiber-apis/databases"
	"fiber-apis/models"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	Username  string `json:"username"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Password  string `json:"password"`
}

func GetUserResponse(user models.User) User {
	return User{Username: user.Username, Firstname: user.Firstname, Lastname: user.Lastname, Password: user.Password}
}

func RegisterHandler(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(404).JSON("Bad user info")
	}
	databases.DataBase.Create(&user)
	userResponse := GetUserResponse(user)

	return c.Status(200).JSON(userResponse)
}

func ViewUsersHandler(c *fiber.Ctx) error {
	var users []models.User

	databases.DataBase.Find(&users)
	var usersResponse []User
	for _, user := range users {
		userResponse := GetUserResponse(user)
		usersResponse = append(usersResponse, userResponse)
	}
	return c.Status(200).JSON(usersResponse)
}
