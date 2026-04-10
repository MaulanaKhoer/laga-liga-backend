package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"laga-liga-backend/config"
	"laga-liga-backend/models"
)

// GetTournamentTeams mengambil semua tim yang terdaftar dalam sebuah turnamen.
// Endpoint: GET /api/tournaments/:id/teams
func GetTournamentTeams(c *gin.Context) {
	tournamentID := c.Param("id")

	// Pastikan turnamen ada dulu
	var tournament models.Tournament
	if err := config.DB.First(&tournament, tournamentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Turnamen tidak ditemukan!"})
		return
	}

	// Preload Teams: GORM akan JOIN ke tabel tournament_teams → teams
	// dan ikutkan pemain masing-masing tim (Preload bersarang)
	config.DB.Preload("Teams.Players").First(&tournament, tournamentID)

	c.JSON(http.StatusOK, gin.H{
		"tournament": gin.H{
			"id":   tournament.ID,
			"name": tournament.Name,
		},
		"total_teams": len(tournament.Teams),
		"teams":       tournament.Teams,
	})
}

// RegisterTeamToTournament mendaftarkan sebuah tim ke dalam turnamen.
// Endpoint: POST /api/tournaments/:id/teams
// Body JSON: { "team_id": 2 }
// Hanya admin yang bisa akses (dijaga di routes).
func RegisterTeamToTournament(c *gin.Context) {
	tournamentID := c.Param("id")

	// ── STEP 1: Ambil turnamen ─────────────────────────────────────────────
	var tournament models.Tournament
	if err := config.DB.Preload("Teams").First(&tournament, tournamentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Turnamen tidak ditemukan!"})
		return
	}

	// ── STEP 2: Parse body request ─────────────────────────────────────────
	var input struct {
		TeamID uint `json:"team_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field team_id wajib diisi!"})
		return
	}

	// ── STEP 3: Pastikan tim yang mau didaftarkan ada ──────────────────────
	var team models.Team
	if err := config.DB.First(&team, input.TeamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tim tidak ditemukan!"})
		return
	}

	// ── STEP 4: Cek apakah tim sudah terdaftar di turnamen ini ────────────
	for _, t := range tournament.Teams {
		if t.ID == team.ID {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Tim sudah terdaftar di turnamen ini!",
			})
			return
		}
	}

	// ── STEP 5: Cek batas maksimal tim ────────────────────────────────────
	if len(tournament.Teams) >= tournament.MaxTeams {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Turnamen sudah penuh! Kapasitas maksimal tercapai.",
			"max_teams": tournament.MaxTeams,
		})
		return
	}

	// ── STEP 6: Daftarkan tim ke turnamen (Many2Many Association) ─────────
	// GORM Append akan INSERT satu baris ke tabel pivot tournament_teams
	if err := config.DB.Model(&tournament).Association("Teams").Append(&team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan tim ke turnamen"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tim berhasil didaftarkan ke turnamen!",
		"tournament": gin.H{
			"id":   tournament.ID,
			"name": tournament.Name,
		},
		"team": gin.H{
			"id":   team.ID,
			"name": team.Name,
			"city": team.City,
		},
		"total_teams_terdaftar": len(tournament.Teams) + 1,
	})
}

// RemoveTeamFromTournament mengeluarkan tim dari sebuah turnamen.
// Endpoint: DELETE /api/tournaments/:id/teams/:team_id
// Hanya admin yang bisa akses (dijaga di routes).
func RemoveTeamFromTournament(c *gin.Context) {
	tournamentID := c.Param("id")
	teamID := c.Param("team_id")

	// ── STEP 1: Ambil turnamen ─────────────────────────────────────────────
	var tournament models.Tournament
	if err := config.DB.First(&tournament, tournamentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Turnamen tidak ditemukan!"})
		return
	}

	// ── STEP 2: Ambil tim ─────────────────────────────────────────────────
	var team models.Team
	if err := config.DB.First(&team, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tim tidak ditemukan!"})
		return
	}

	// ── STEP 3: Hapus relasi (DELETE dari tabel pivot tournament_teams) ────
	// GORM Delete Association hanya menghapus baris di tabel pivot,
	// bukan menghapus data tim itu sendiri dari tabel teams.
	if err := config.DB.Model(&tournament).Association("Teams").Delete(&team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengeluarkan tim dari turnamen"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tim berhasil dikeluarkan dari turnamen!",
		"tournament": gin.H{
			"id":   tournament.ID,
			"name": tournament.Name,
		},
		"team": gin.H{
			"id":   team.ID,
			"name": team.Name,
		},
	})
}
