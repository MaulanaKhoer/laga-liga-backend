package routes

import (
	"github.com/gofiber/fiber/v2"
	"laga-liga-backend/internal/handlers"
	"laga-liga-backend/internal/middleware"
)

func setupMatchEventRoutes(protected fiber.Router) {
	matchEvents := protected.Group("/matches/:id/events")
	
	matchEvents.Get("", handlers.GetMatchEvents)
	matchEvents.Post("", middleware.RequireRole("admin"), handlers.AddMatchEvent)
	matchEvents.Delete("/:event_id", middleware.RequireRole("admin"), handlers.DeleteMatchEvent)
}
