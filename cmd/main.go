// cmd/main.go
package main

import (
	"github.com/diana-yelemes/habit_tracker/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
)

func main() {
	database.ConnectDb()
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})
	app.Use(cors.New())

	setupRoutes(app)
	app.Static("/", "./public")

	app.Listen(":3000")
}
