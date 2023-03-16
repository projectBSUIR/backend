package routes

import (
	"fiber-apis/models"
	"github.com/gofiber/fiber/v2"
	"log"
)

type UserResponse struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email, omitempty""`
}

func (user *UserResponse) GetUserModel() models.User {
	return models.User{ID: 0, Login: user.Login, Password: user.Password, Email: user.Email, Status: models.UserStatus(models.Participant)}
}

func RegisterHandler(c *fiber.Ctx) error {
	var user UserResponse

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON("")
	}
	userModel := user.GetUserModel()
	err := userModel.Register()

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	///generate JWT
	return c.Status(200).JSON("Register Successfull")
}

func LoginUserHandler(c *fiber.Ctx) error {
	var user UserResponse

	if err := c.BodyParser(&user); err != nil {
		log.Fatal(string(c.Body()))
		return c.Status(400).JSON("")
	}
	userModel := user.GetUserModel()
	err := userModel.LogIn()

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	///generate JWT
	return c.Status(200).JSON("Login successful")
}
