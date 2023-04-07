package middlewares

import (
	"fiber-apis/models"
	"fiber-apis/token"
	"github.com/gofiber/fiber/v2"
)

func Participant(c *fiber.Ctx) error {
	signedAccessToken := c.GetReqHeaders()["Authorization"]
	refreshToken := c.Cookies("refresh_token")

	if refreshToken == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	accessClaims, err := token.ValidateToken(signedAccessToken)
	if err != nil {
		if err.Error() == token.JWTErrTokenExpired.Error() {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	refreshClaims, err := token.GetClaims(refreshToken)

	if refreshClaims.User != accessClaims.User {
		return c.SendStatus(fiber.StatusForbidden)
	}

	return c.Next()
}

func Coach(c *fiber.Ctx) error {
	signedAccessToken := c.GetReqHeaders()["Authorization"]
	refreshToken := c.Cookies("refresh_token")

	if refreshToken == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	accessClaims, err := token.ValidateToken(signedAccessToken)
	if err != nil {
		if err.Error() == token.JWTErrTokenExpired.Error() {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	refreshClaims, err := token.GetClaims(refreshToken)

	if refreshClaims.User != accessClaims.User {
		return c.SendStatus(fiber.StatusForbidden)
	}

	if accessClaims.Status < models.Admin {
		return c.SendStatus(fiber.StatusForbidden)
	}
	return c.Next()
}

func Admin(c *fiber.Ctx) error {
	signedAccessToken := c.GetReqHeaders()["Authorization"]
	refreshToken := c.Cookies("refresh_token")

	if refreshToken == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	accessClaims, err := token.ValidateToken(signedAccessToken)
	if err != nil {
		if err.Error() == token.JWTErrTokenExpired.Error() {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	refreshClaims, err := token.GetClaims(refreshToken)

	if refreshClaims.User != accessClaims.User {
		return c.SendStatus(fiber.StatusForbidden)
	}
	if accessClaims.Status != models.Admin {
		return c.SendStatus(fiber.StatusForbidden)
	}

	return c.Next()
}
