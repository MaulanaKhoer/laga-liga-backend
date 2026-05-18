package routes

import (
	"github.com/gofiber/fiber/v2"
	"laga-liga-backend/internal/handlers"
)

func setupAuthRoutes(api fiber.Router) {
	auth := api.Group("/auth")
	auth.Post("/register", handlers.Register)
	auth.Post("/login", handlers.Login)
}
