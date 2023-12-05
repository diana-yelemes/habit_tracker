// cmd/routes.go

package main

import (
	"github.com/diana-yelemes/habit_tracker/handlers"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", handlers.Home)
	app.Get("/habit", handlers.GetAllUserHabits)
	app.Post("/habit", handlers.CreateHabit)
	app.Put("/habit/:id", handlers.UpdateHabit)
	app.Delete("/habit/:id", handlers.DeleteHabit)
	app.Get("/habit/:id", handlers.GetHabitByID)
	app.Get("/habits/filter", handlers.FilterHabits)
	app.Put("/habit/complete/:id", handlers.CompleteHabit)
	app.Put("/habit/undo-complete/:id", handlers.UndoCompleteHabit)
	app.Get("/habits/completed", handlers.GetCompletedHabits)
	app.Get("/habits/incomplete", handlers.GetIncompleteHabits)
	app.Get("/habits/statistics", handlers.GetHabitStatistics)
}
