package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupCORS mengembalikan handler untuk menangani Cross-Origin Resource Sharing.
// Ini penting agar Frontend (React/Vue/Next.js) bisa mengakses API ini.
func SetupCORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		// Izinkan akses dari origin mana saja (untuk development)
		// Jika production, sebaiknya spesifik: []string{"https://lagaliga.id"}
		AllowAllOrigins:  true,
		
		// Method HTTP yang diizinkan
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		
		// Header yang diizinkan dikirim oleh client
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		
		// Header yang diizinkan dibaca oleh client dari response
		ExposeHeaders:    []string{"Content-Length"},
		
		// Izinkan pengiriman cookies/auth headers
		AllowCredentials: true,
		
		// Berapa lama hasil pre-flight request (OPTIONS) disimpan di browser
		MaxAge: 12 * time.Hour,
	})
}
