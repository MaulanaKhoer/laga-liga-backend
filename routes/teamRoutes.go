package routes

import (
	"github.com/gin-gonic/gin"
	"laga-liga-backend/controllers"
	"laga-liga-backend/middleware"
)

// setupTeamRoutes mendaftarkan semua route untuk modul Tim.
func setupTeamRoutes(protected *gin.RouterGroup) {
	// GET: semua user login bisa lihat
	protected.GET("/teams", controllers.GetTeams)
	protected.GET("/teams/:id", controllers.GetTeamByID)

	// POST, PUT, DELETE: hanya admin
	protected.POST("/teams", middleware.RequireRole("admin"), controllers.CreateTeam)
	protected.PUT("/teams/:id", middleware.RequireRole("admin"), controllers.UpdateTeam)
	protected.DELETE("/teams/:id", middleware.RequireRole("admin"), controllers.DeleteTeam)

	// Tambah pemain ke tim — nested route
	// Contoh: POST /api/teams/1/players
	protected.POST("/teams/:id/players", middleware.RequireRole("admin"), controllers.AddPlayerToTeam)
}
