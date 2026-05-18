package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"laga-liga-backend/internal/config"
	"laga-liga-backend/internal/middleware"
	"laga-liga-backend/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: file .env tidak ditemukan, menggunakan nilai default OS.")
	}

	config.ConnectDatabase()

	app := fiber.New()

	app.Use(middleware.SetupCORS())

	routes.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("🚀 Server berjalan di http://localhost:%s", port)
	app.Listen(":" + port)
}
