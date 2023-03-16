package main

import (
	"fiber-apis/databases"
)

func main() {
	//app := fiber.New()

	err := databases.ConnectDB()
	if err != nil {
		panic(err.Error())
	}

	//routes.Setup(app)

	//app.Listen(":5000")
}
