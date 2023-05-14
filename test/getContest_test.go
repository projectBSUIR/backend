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

func TestGetContests(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	app := fiber.New(fiber.Config{
		AppName:   "TestServer",
		BodyLimit: 128 * 1024 * 1024,
	})

	err := databases.ConnectDB()
	if err != nil {
		panic(err.Error())
	}

	routes.Setup(app)

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

	token := string(body)

	var unmarshaled Token

	err = json.Unmarshal([]byte(token), &unmarshaled)
	if err != nil {
		t.Fatal("Error decoding JSON:", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contest := httptest.NewRequest(http.MethodGet, "/ownContests", nil)
			contest.Header.Set("Authorization", unmarshaled.Token)

			resp, err := app.Test(contest)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != http.StatusInternalServerError {
				t.Errorf("expected status %d; got %d", http.StatusOK, resp.StatusCode)
			}
		})
	}
}
