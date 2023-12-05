// cmd/routes.go

package main

import (
	"github.com/diana-yelemes/habit_tracker/handlers"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", handlers.Home)
	app.Get("/habit", func(c *fiber.Ctx) error { return handlers.GetAllUserHabits(c, setup.DB) })
	app.Post("/habit", func(c *fiber.Ctx) error { return handlers.CreateHabit(c, setup.DB) })
	app.Post("/habitlist", func(c *fiber.Ctx) error { return handlers.CreateHabitList(c, setup.DB) })
	app.Put("/habit", func(c *fiber.Ctx) error { return handlers.UpdateHabit(c, setup.DB) })
	app.Put("/habitlist", func(c *fiber.Ctx) error { return handlers.UpdateHabitList(c, setup.DB) })
	app.Delete("/habit", func(c *fiber.Ctx) error { return handlers.DeleteHabit(c, setup.DB) })
	app.Delete("/habitlist", func(c *fiber.Ctx) error { return handlers.DeleteHabitList(c, setup.DB) })
}
