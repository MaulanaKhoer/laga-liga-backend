package routes

import (
	"github.com/gin-gonic/gin"
	"laga-liga-backend/controllers"
	"laga-liga-backend/middleware"
)

// SetupRoutes adalah pintu masuk utama semua route.
// File ini hanya mengatur struktur besar — detail tiap modul ada di file masing-masing.
func SetupRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong dari laga-liga backend (Golang)!"})
	})

	api := r.Group("/api")

	// ── PUBLIC ROUTES ──────────────────────────────────────────────────────
	setupAuthRoutes(api)

	// ── PROTECTED ROUTES (wajib login) ────────────────────────────────────
	protected := api.Group("/")
	protected.Use(middleware.RequireAuth)
	{
		protected.GET("/me", controllers.GetMe)

		// Daftarkan routes per modul — tambah modul baru cukup 1 baris di sini
		setupTournamentRoutes(protected)
		// setupTeamRoutes(protected)   ← nanti kalau buat modul Team
		// setupMatchRoutes(protected)  ← nanti kalau buat modul Match
	}
}
