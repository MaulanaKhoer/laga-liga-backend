package routes

import (
	"github.com/gofiber/fiber/v2"
	"laga-liga-backend/internal/handlers"
	"laga-liga-backend/internal/middleware"
)

func setupTournamentRoutes(protected fiber.Router) {
	protected.Get("/tournaments", handlers.GetTournaments)
	protected.Get("/tournaments/:id", handlers.GetTournamentByID)
	protected.Get("/sports", handlers.GetSports)

	protected.Post("/tournaments", middleware.RequireRole("admin"), handlers.CreateTournament)
	protected.Put("/tournaments/:id", middleware.RequireRole("admin"), handlers.UpdateTournament)
	protected.Delete("/tournaments/:id", middleware.RequireRole("admin"), handlers.DeleteTournament)

	protected.Get("/tournaments/:id/teams", handlers.GetTournamentTeams)
	protected.Post("/tournaments/:id/teams", middleware.RequireRole("admin"), handlers.RegisterTeamToTournament)
	protected.Delete("/tournaments/:id/teams/:team_id", middleware.RequireRole("admin"), handlers.RemoveTeamFromTournament)

	protected.Get("/tournaments/:id/standings", handlers.GetTournamentStandings)
	protected.Get("/tournaments/:id/top-scorers", handlers.GetTopScorers)
}
