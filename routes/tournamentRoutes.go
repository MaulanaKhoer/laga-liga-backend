package routes

import (
	"github.com/gin-gonic/gin"
	"laga-liga-backend/controllers"
	"laga-liga-backend/middleware"
)

// setupTournamentRoutes mendaftarkan semua route turnamen.
// Menerima `protected` group yang sudah terpasang middleware RequireAuth.
func setupTournamentRoutes(protected *gin.RouterGroup) {
	// GET: semua user yang login boleh akses
	protected.GET("/tournaments", controllers.GetTournaments)
	protected.GET("/tournaments/:id", controllers.GetTournamentByID)

	// POST, PUT, DELETE: hanya admin
	protected.POST("/tournaments", middleware.RequireRole("admin"), controllers.CreateTournament)
	protected.PUT("/tournaments/:id", middleware.RequireRole("admin"), controllers.UpdateTournament)
	protected.DELETE("/tournaments/:id", middleware.RequireRole("admin"), controllers.DeleteTournament)

	// ── Manajemen Tim dalam Turnamen (pivot tournament_teams) ──────────────
	// GET  /api/tournaments/:id/teams         → lihat semua tim yang ikut
	// POST /api/tournaments/:id/teams         → daftarkan tim ke turnamen (admin)
	// DELETE /api/tournaments/:id/teams/:team_id → keluarkan tim dari turnamen (admin)
	protected.GET("/tournaments/:id/teams", controllers.GetTournamentTeams)
	protected.POST("/tournaments/:id/teams", middleware.RequireRole("admin"), controllers.RegisterTeamToTournament)
	protected.DELETE("/tournaments/:id/teams/:team_id", middleware.RequireRole("admin"), controllers.RemoveTeamFromTournament)

	// ── Klasemen (Standings) ────────────────────────────────────────────────
	protected.GET("/tournaments/:id/standings", controllers.GetTournamentStandings)
	protected.GET("/tournaments/:id/top-scorers", controllers.GetTopScorers)
}
