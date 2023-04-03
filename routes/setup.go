package routes

import (
	"fiber-apis/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/register", RegisterHandler)
	app.Post("/login", LoginUserHandler)
	app.Post("/logout", Logout)
	app.Get("/refresh", RefreshToken)

	app.Use(middlewares.Participant)
	app.Get("/check", CheckHandler)

	app.Use(middlewares.Admin)
	app.Post("/addProblem", AddProblem)
	app.Post("/createContest", CreateContest)
}
