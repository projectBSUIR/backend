package test

import (
	"fiber-apis/databases"
	"fiber-apis/routes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRegisterHandler(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		json    string
		args    args
		wantErr bool
	}{
		{
			name: "register user 1",
			json: `
			{
				"login": "testuser1",
				"password": "test1",
				"email": "email1"
		    }
			`,
			args: args{
				c: &fiber.Ctx{},
			},
			wantErr: true,
		},
		{
			name: "register user 2",
			json: `
			{
				"login": "testuser2",
				"password": "test2",
				"email": "email2"
		    }
			`,
			args: args{
				c: &fiber.Ctx{},
			},
			wantErr: true,
		},
		{
			name: "register user 3",
			json: `
			{
				"login": "testuser3",
				"password": "test3",
				"email": "email3"
		    }
			`,
			args: args{
				c: &fiber.Ctx{},
			},
			wantErr: true,
		},
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(tt.json))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status %d; got %d", http.StatusOK, resp.StatusCode)
			}
		})
	}
}
