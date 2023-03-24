package routes

import (
	"fiber-apis/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/register", RegisterHandler)
	app.Post("/login", LoginUserHandler)
	app.Get("/refresh", RefreshToken)
	app.Post("/logout", Logout)
	app.Post("/addContest", ContestHandler)

	app.Use(middlewares.Participant)
	app.Get("/check", CheckHandler)
}
