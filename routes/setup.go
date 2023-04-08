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
	app.Get("/contests", ViewContests)
	app.Get("/contest/:contestId", ViewProblems)

	app.Use(middlewares.Participant)
	app.Get("/check", CheckHandler)
	app.Post("/submit", SubmitSolution)

	app.Use(middlewares.Admin)
	app.Post("/addProblem", AddProblem)
	app.Post("/createContest", CreateContest)
}
