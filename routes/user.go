package routes

import (
	"fiber-apis/models"
	"fiber-apis/token"
	"github.com/gofiber/fiber/v2"
	"time"
)

type UserResponse struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email, omitempty""`
}

func (user *UserResponse) GetUserModel() models.User {
	return models.User{ID: 0, Login: user.Login, Password: user.Password, Email: user.Email, Status: models.Participant}
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

	return Authentificate(c, &userModel)
}

func LoginUserHandler(c *fiber.Ctx) error {
	var user UserResponse

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	userModel := user.GetUserModel()
	err := userModel.LogIn()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	return Authentificate(c, &userModel)
}

func Authentificate(c *fiber.Ctx, userModel *models.User) error {
	accessToken, err := token.GenerateAccessToken(c, token.GetJWTClaim(userModel, time.Now().Add(time.Minute*5)))
	if err != nil {
		return err
	}
	err = token.GenerateRefreshToken(c, token.GetJWTClaim(userModel, time.Now().Add(time.Hour*72)))
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": accessToken,
	})
}

func RefreshToken(c *fiber.Ctx) error {
	signedToken := c.Cookies("refresh_token")
	if signedToken == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	accessToken, err := token.Refresh(c, signedToken)

	if err != nil {
		if err.Error() == token.JWTErrTokenExpired.Error() {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": accessToken,
	})
}

func Logout(c *fiber.Ctx) error {
	signedToken := c.Cookies("refresh_token")
	if signedToken == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	token.SetRefreshTokenCookie(c, signedToken, time.Now())
	return c.SendStatus(fiber.StatusOK)
}
