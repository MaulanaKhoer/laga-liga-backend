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
}
