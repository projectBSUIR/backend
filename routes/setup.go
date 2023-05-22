package routes

import (
	"fiber-apis/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Setup(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Accept, Content-Type, Content-Length, Accept-Encoding, Authorization",
		AllowOrigins:     "http://*",
		AllowCredentials: true,
	}))

	app.Post("/login", LoginUserHandler)
	app.Post("/logout", Logout)
	app.Post("/register", RegisterHandler)
	app.Get("/contests", ViewContests)
	app.Get("/contest/:contestId", ViewProblems)
	app.Get("/refresh", RefreshToken)
	app.Get("/contest/:contestId/standings", GetResultsTable)

	testingMachineApp := app.Group("/testing", middlewares.TestMachine)
	testingMachineApp.Get("/extractSubmission", ExtractSubmissionFromTestingQueue)
	testingMachineApp.Get("/extractTestingFiles", ExtractFilesForTesting)
	testingMachineApp.Get("/extractProblemTests", ExtractProblemTests)
	testingMachineApp.Post("/setVerdict", SetVerdict)

	app.Use(middlewares.Participant)
	app.Get("/submissions/:problemId", GetSubmissions)
	app.Get("/contest/:contestId/submissions", GetAllSubmissions)
	app.Post("/submit", SubmitSolution)
	app.Get("/ownContests", GetContests)

	app.Use(middlewares.Coach)
	app.Post("/createContest", CreateContest)
	app.Post("/addProblem", AddProblem)

	app.Use(middlewares.Admin)
	app.Post("/setCoach", SetCoach)
}
