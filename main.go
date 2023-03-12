package main

import (
	"fiber-apis/databases"
	"fiber-apis/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	databases.Connect(app)
	routes.Setup(app)

	app.Listen(":5000")
}
