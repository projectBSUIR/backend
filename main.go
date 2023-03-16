package main

import (
	"fiber-apis/databases"
	"fiber-apis/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	err := databases.ConnectDB()
	if err != nil {
		panic(err.Error())
	}

	routes.Setup(app)

	app.Listen(":5000")
}
