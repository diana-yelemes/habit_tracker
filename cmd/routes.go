// cmd/routes.go

package main

import (
	"github.com/diana-yelemes/habit_tracker/handlers"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", handlers.GetAllUserHabits)
	app.Post("/habit", handlers.CreateHabit)
	app.Put("/habit/:id", handlers.UpdateHabit)
	app.Delete("/habitdelete/:id", handlers.DeleteHabit)
	app.Get("/habit/:id", handlers.GetHabitByID)
	app.Get("/habits/filter", handlers.FilterHabits)
	app.Put("/habit/complete/:id", handlers.CompleteHabit)
	app.Put("/habit/undo-complete/:id", handlers.UndoCompleteHabit)
	app.Get("/habits/completed", handlers.GetCompletedHabits)
	app.Get("/habits/incomplete", handlers.GetIncompleteHabits)
	app.Get("/habits/statistics", handlers.GetHabitStatistics)
	app.Get("/habit", handlers.NewHabitView)
	app.Get("/delete-habit", handlers.DeleteHabitView)
	app.Put("/habit/updateRepeatCount/:id", handlers.UpdateRepeatCount)
}
