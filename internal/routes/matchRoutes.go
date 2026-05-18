package routes

import (
	"github.com/gofiber/fiber/v2"
	"laga-liga-backend/internal/handlers"
	"laga-liga-backend/internal/middleware"
)

func setupMatchRoutes(protected fiber.Router) {
	protected.Get("/matches", handlers.GetMatches)
	protected.Get("/matches/:id", handlers.GetMatchByID)

	protected.Post("/matches", middleware.RequireRole("admin"), handlers.CreateMatch)
	protected.Put("/matches/:id", middleware.RequireRole("admin"), handlers.UpdateMatch)
	protected.Delete("/matches/:id", middleware.RequireRole("admin"), handlers.DeleteMatch)

	protected.Patch("/matches/:id/score", middleware.RequireRole("admin"), handlers.UpdateScore)

	protected.Get("/tournaments/:id/matches", handlers.GetMatchesByTournament)
}
