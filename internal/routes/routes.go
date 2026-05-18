package routes

import (
	"github.com/gofiber/fiber/v2"

	"laga-liga-backend/internal/handlers"
	"laga-liga-backend/internal/middleware"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "pong dari laga-liga backend (Golang Fiber)!"})
	})

	api := app.Group("/api")

	setupAuthRoutes(api)

	protected := api.Group("/")
	protected.Use(middleware.RequireAuth)

	protected.Get("/me", handlers.GetMe)

	setupTournamentRoutes(protected)
	setupTeamRoutes(protected)
	setupMatchRoutes(protected)
	setupMatchEventRoutes(protected)
}
