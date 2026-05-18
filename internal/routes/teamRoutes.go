package routes

import (
	"github.com/gofiber/fiber/v2"
	"laga-liga-backend/internal/handlers"
	"laga-liga-backend/internal/middleware"
)

func setupTeamRoutes(protected fiber.Router) {
	protected.Get("/teams", handlers.GetTeams)
	protected.Get("/teams/:id", handlers.GetTeamByID)

	protected.Post("/teams", middleware.RequireRole("admin"), handlers.CreateTeam)
	protected.Put("/teams/:id", middleware.RequireRole("admin"), handlers.UpdateTeam)
	protected.Delete("/teams/:id", middleware.RequireRole("admin"), handlers.DeleteTeam)

	protected.Post("/teams/:id/players", middleware.RequireRole("admin"), handlers.AddPlayerToTeam)
}
