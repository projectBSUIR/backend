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

func TestLoginUserHandler(t *testing.T) {
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
			name: "admin",
			json: `
			{
				"login": "admin",
				"password": "22df674c179820a70cbbe183be510f2c781edc3b5286c7f40f97fc6c8ee75101"
			}
			`,
			args: args{
				c: &fiber.Ctx{},
			},
			wantErr: true,
		},
		{
			name: "testmachine",
			json: `
			{
				"login": "testmachine",
				"password": "202979ff6105a2c7859b95efb411adc7392e43b6aae92a88d8992aba90fe83f4"
			}
			`,
			args: args{
				c: &fiber.Ctx{},
			},
			wantErr: true,
		},
		{
			name: "testuser",
			json: `
			{
				"login": "testuser",
				"password": "test1234"
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
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.json))
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
