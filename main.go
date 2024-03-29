package main

import (
	"fiber-apis/databases"
	"fiber-apis/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"os"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName:   "BetterSolve",
		BodyLimit: 128 * 1024 * 1024,
	})

	file, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	err = databases.ConnectDB()
	if err != nil {
		panic(err.Error())
	}
	defer databases.DataBase.Close()

	app.Use(logger.New())
	routes.Setup(app)

	app.Listen(":5000")
}
