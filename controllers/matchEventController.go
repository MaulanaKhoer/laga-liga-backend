package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"laga-liga-backend/config"
	"laga-liga-backend/models"
)

// AddMatchEvent menambahkan kejadian (gol/kartu) ke dalam pertandingan. Hanya admin.
// Endpoint: POST /api/matches/:id/events
func AddMatchEvent(c *gin.Context) {
	matchID := c.Param("id")

	// 1. Pastikan pertandingan ada
	var match models.Match
	if err := config.DB.First(&match, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertandingan tidak ditemukan!"})
		return
	}

	// 2. Bind input
	var input struct {
		PlayerID uint   `json:"player_id" binding:"required"`
		Type     string `json:"type"      binding:"required"` // goal, yellow_card, red_card, own_goal
		Minute   int    `json:"minute"    binding:"required"`
		Note     string `json:"note"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap: " + err.Error()})
		return
	}

	// 3. Pastikan pemain ada
	var player models.Player
	if err := config.DB.First(&player, input.PlayerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pemain tidak ditemukan!"})
		return
	}

	// 4. Pastikan pemain tersebut membela salah satu tim yang sedang bertanding
	if player.TeamID != match.HomeTeamID && player.TeamID != match.AwayTeamID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Pemain ini tidak terdaftar di tim yang sedang bertanding!",
		})
		return
	}

	// 5. Simpan event
	event := models.MatchEvent{
		MatchID:  match.ID,
		PlayerID: player.ID,
		TeamID:   player.TeamID,
		Type:     input.Type,
		Minute:   input.Minute,
		Note:     input.Note,
	}

	if err := config.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan event"})
		return
	}

	// Load ulang dengan detail pemain agar response cantik
	config.DB.Preload("Player").Preload("Team").First(&event, event.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Event berhasil dicatat!",
		"data":    event,
	})
}

// GetMatchEvents mengambil semua kejadian dalam satu pertandingan.
// Endpoint: GET /api/matches/:id/events
func GetMatchEvents(c *gin.Context) {
	matchID := c.Param("id")

	var events []models.MatchEvent
	config.DB.
		Preload("Player").
		Preload("Team").
		Where("match_id = ?", matchID).
		Order("minute ASC").
		Find(&events)

	c.JSON(http.StatusOK, gin.H{"data": events})
}

// DeleteMatchEvent menghapus kejadian pertandingan. Hanya admin.
// Endpoint: DELETE /api/matches/:id/events/:event_id
func DeleteMatchEvent(c *gin.Context) {
	eventID := c.Param("event_id")

	var event models.MatchEvent
	if err := config.DB.First(&event, eventID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event tidak ditemukan!"})
		return
	}

	config.DB.Delete(&event)
	c.JSON(http.StatusOK, gin.H{"message": "Event berhasil dihapus!"})
}
