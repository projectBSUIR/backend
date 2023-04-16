package routes

import (
	"fiber-apis/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/login", LoginUserHandler)
	app.Post("/logout", Logout)
	app.Post("/register", RegisterHandler)
	app.Get("/contests", ViewContests)
	app.Get("/contest/:contestId", ViewProblems)
	app.Get("/refresh", RefreshToken)
	app.Get("/standings", GetResultsTable)

	app.Use(middlewares.Participant)
	app.Get("/check", CheckHandler)
	app.Post("/submit", SubmitSolution)
	app.Get("/ownContests", GetContests)

	app.Use(middlewares.Coach)
	app.Post("/createContest", CreateContest)
	app.Post("/addProblem", AddProblem)

	app.Use(middlewares.Admin)
	app.Post("/setCoach", SetCoach)
}
