package routes

import (
	"github.com/gin-gonic/gin"
	"laga-liga-backend/controllers"
	"laga-liga-backend/middleware"
)

// setupMatchEventRoutes mendaftarkan route untuk kejadian pertandingan (gol/kartu).
func setupMatchEventRoutes(protected *gin.RouterGroup) {
	// Di bawah /matches/:id/events
	matchEvents := protected.Group("/matches/:id/events")
	{
		matchEvents.GET("", controllers.GetMatchEvents)
		matchEvents.POST("", middleware.RequireRole("admin"), controllers.AddMatchEvent)
		matchEvents.DELETE("/:event_id", middleware.RequireRole("admin"), controllers.DeleteMatchEvent)
	}
}
