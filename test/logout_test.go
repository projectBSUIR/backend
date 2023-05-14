package test

import (
	"encoding/json"
	"fiber-apis/databases"
	"fiber-apis/routes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type Token struct {
	Token string `json:"access_token"`
}

func TestLogout(t *testing.T) {
	app := fiber.New(fiber.Config{
		AppName:   "TestServer",
		BodyLimit: 128 * 1024 * 1024,
	})

	err := databases.ConnectDB()
	if err != nil {
		panic(err.Error())
	}

	routes.Setup(app)

	token := ""

	t.Run("test logout", func(t *testing.T) {
		log := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`
		{
			"login": "admin",
			"password": "22df674c179820a70cbbe183be510f2c781edc3b5286c7f40f97fc6c8ee75101"
		}
		`))
		log.Header.Set("Content-Type", "application/json")

		loginResponse, err := app.Test(log)
		if err != nil {
			t.Fatal(err)
		}

		if loginResponse.StatusCode != http.StatusOK {
			t.Errorf("expected status %d; got %d", http.StatusOK, loginResponse.StatusCode)
		}

		body, err := ioutil.ReadAll(loginResponse.Body)
		if err != nil {
			t.Fatal(err)
		}

		token = string(body)

		var unmarshaled Token

		err = json.Unmarshal([]byte(token), &unmarshaled)
		if err != nil {
			t.Fatal("Error decoding JSON:", err)
		}

		logout := httptest.NewRequest(http.MethodPost, "/logout", nil)
		logout.Header.Set("Authorization", unmarshaled.Token)

		logoutResponse, logoutErr := app.Test(logout)
		if logoutErr != nil {
			t.Fatal(err)
		}

		if logoutResponse.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected status %d; got %d", http.StatusOK, logoutResponse.StatusCode)
		}
	})
}
