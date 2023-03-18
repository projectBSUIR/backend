package token

import (
	"encoding/json"
	"errors"
	"fiber-apis/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var JWTErrTokenExpired = errors.New("Token is expired")

var jwtKey = []byte("secret")

type JWTClaim struct {
	User   string            `json:"user"`
	Status models.UserStatus `json:"status"`
	jwt.RegisteredClaims
}

func SetRefreshTokenCookie(c *fiber.Ctx, refreshToken string, expiresAt time.Time) {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  expiresAt,
		HTTPOnly: true,
	})
}

func GetJWTClaim(model *models.User, expirationTime time.Time) JWTClaim {
	return JWTClaim{
		User:   model.Login,
		Status: model.Status,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
}

func GenerateAccessToken(c *fiber.Ctx, claims JWTClaim) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", c.SendStatus(fiber.StatusInternalServerError)
	}
	return accessToken, nil
}

func GenerateRefreshToken(c *fiber.Ctx, claims JWTClaim) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshToken, err := token.SignedString(jwtKey)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	SetRefreshTokenCookie(c, refreshToken, claims.ExpiresAt.Time)
	return nil
}

func ValidateToken(signedToken string) (*JWTClaim, error) {
	claims, err := GetClaims(signedToken)

	if err != nil {
		return &JWTClaim{}, err
	}
	return claims, nil
}

func GetClaims(signedToken string) (*JWTClaim, error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	claims := &JWTClaim{}
	if err != nil {
		return claims, err
	}
	if !token.Valid {
		return claims, errors.New("Token is not valid")
	}
	mapClaims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return claims, errors.New("Invalid token")
	}

	jsonClaims, err := json.Marshal(mapClaims)

	if err != nil {
		return claims, err
	}
	json.Unmarshal(jsonClaims, &claims)
	return claims, err
}

func Refresh(c *fiber.Ctx, signedToken string) (string, error) {
	claims, err := GetClaims(signedToken)
	if err != nil {
		return "", err
	}
	signedAccessToken := c.GetReqHeaders()["Authorization"]
	if signedAccessToken == "" {
		return "", errors.New("Headers haven't an access token")
	}
	claimsAccessToken, err := ValidateToken(signedAccessToken)

	if err != nil && err.Error() != JWTErrTokenExpired.Error() {
		return "", err
	}
	if err == nil {
		if claims.User != claimsAccessToken.User {
			return "", errors.New("User from access token and user from refresh token are different")
		}
		return signedAccessToken, nil
	}
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Second * 20))

	return GenerateAccessToken(c, *claims)
}
