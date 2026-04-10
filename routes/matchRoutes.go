package routes

import (
	"github.com/gin-gonic/gin"
	"laga-liga-backend/controllers"
	"laga-liga-backend/middleware"
)

// setupMatchRoutes mendaftarkan semua route untuk modul Pertandingan.
//
// Struktur endpoint:
//
//	GET    /api/matches                      → semua pertandingan (bisa filter ?tournament_id=&status=)
//	GET    /api/matches/:id                  → 1 pertandingan
//	POST   /api/matches                      → buat jadwal baru (admin)
//	PUT    /api/matches/:id                  → update detail jadwal (admin)
//	PATCH  /api/matches/:id/score            → update skor & status (admin)
//	DELETE /api/matches/:id                  → hapus pertandingan (admin)
//
// Endpoint nested (di bawah /tournaments):
//
//	GET    /api/tournaments/:id/matches      → pertandingan dalam 1 turnamen
func setupMatchRoutes(protected *gin.RouterGroup) {
	// ── Match global ──────────────────────────────────────────────────────
	protected.GET("/matches", controllers.GetMatches)
	protected.GET("/matches/:id", controllers.GetMatchByID)

	protected.POST("/matches", middleware.RequireRole("admin"), controllers.CreateMatch)
	protected.PUT("/matches/:id", middleware.RequireRole("admin"), controllers.UpdateMatch)
	protected.DELETE("/matches/:id", middleware.RequireRole("admin"), controllers.DeleteMatch)

	// Endpoint khusus update skor — PATCH karena partial update
	protected.PATCH("/matches/:id/score", middleware.RequireRole("admin"), controllers.UpdateScore)

	// ── Match nested dalam tournament ─────────────────────────────────────
	// Catatan: route /tournaments/:id sudah ada di tournamentRoutes.go
	// Kita tambahkan /matches di sini sebagai nested resource
	protected.GET("/tournaments/:id/matches", controllers.GetMatchesByTournament)
}
