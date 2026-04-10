package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"laga-liga-backend/config"
	"laga-liga-backend/models"
)

// GetMatches mengambil semua pertandingan.
// Support query string: ?tournament_id=1  ?status=finished
// Endpoint: GET /api/matches
func GetMatches(c *gin.Context) {
	var matches []models.Match

	query := config.DB.
		Preload("Tournament").
		Preload("HomeTeam").
		Preload("AwayTeam")

	// Filter opsional by tournament
	if tid := c.Query("tournament_id"); tid != "" {
		query = query.Where("tournament_id = ?", tid)
	}

	// Filter opsional by status
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Order("match_date ASC").Find(&matches)

	c.JSON(http.StatusOK, gin.H{
		"total": len(matches),
		"data":  matches,
	})
}

// GetMatchByID mengambil 1 pertandingan berdasarkan ID.
// Endpoint: GET /api/matches/:id
func GetMatchByID(c *gin.Context) {
	var match models.Match
	id := c.Param("id")

	if err := config.DB.
		Preload("Tournament").
		Preload("HomeTeam").
		Preload("AwayTeam").
		First(&match, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertandingan tidak ditemukan!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": match})
}

// GetMatchesByTournament mengambil semua pertandingan dalam 1 turnamen.
// Endpoint: GET /api/tournaments/:id/matches
func GetMatchesByTournament(c *gin.Context) {
	tournamentID := c.Param("id")

	// Pastikan turnamen ada
	var tournament models.Tournament
	if err := config.DB.First(&tournament, tournamentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Turnamen tidak ditemukan!"})
		return
	}

	var matches []models.Match
	config.DB.
		Preload("HomeTeam").
		Preload("AwayTeam").
		Where("tournament_id = ?", tournamentID).
		Order("match_date ASC").
		Find(&matches)

	c.JSON(http.StatusOK, gin.H{
		"tournament": gin.H{"id": tournament.ID, "name": tournament.Name},
		"total":      len(matches),
		"data":       matches,
	})
}

// CreateMatch membuat jadwal pertandingan baru. Hanya admin.
// Endpoint: POST /api/matches
// Body JSON:
//
//	{
//	  "tournament_id": 1,
//	  "home_team_id":  1,
//	  "away_team_id":  2,
//	  "round":         "Grup A",
//	  "venue":         "GBK Jakarta",
//	  "match_date":    "2026-05-10T15:00:00Z"
//	}
func CreateMatch(c *gin.Context) {
	var input struct {
		TournamentID uint   `json:"tournament_id" binding:"required"`
		HomeTeamID   uint   `json:"home_team_id"  binding:"required"`
		AwayTeamID   uint   `json:"away_team_id"  binding:"required"`
		Round        string `json:"round"`
		Venue        string `json:"venue"`
		MatchDate    string `json:"match_date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap: " + err.Error()})
		return
	}

	// Validasi: home dan away tidak boleh sama
	if input.HomeTeamID == input.AwayTeamID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tim kandang dan tandang tidak boleh sama!"})
		return
	}

	// Pastikan turnamen ada
	var tournament models.Tournament
	if err := config.DB.First(&tournament, input.TournamentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Turnamen tidak ditemukan!"})
		return
	}

	// Pastikan kedua tim ada
	var homeTeam, awayTeam models.Team
	if err := config.DB.First(&homeTeam, input.HomeTeamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tim kandang tidak ditemukan!"})
		return
	}
	if err := config.DB.First(&awayTeam, input.AwayTeamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tim tandang tidak ditemukan!"})
		return
	}

	// Parse tanggal pertandingan
	matchDate, err := parseDateString(input.MatchDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Format tanggal tidak valid! Gunakan: 2026-05-10T15:00:00Z",
		})
		return
	}

	match := models.Match{
		TournamentID: input.TournamentID,
		HomeTeamID:   input.HomeTeamID,
		AwayTeamID:   input.AwayTeamID,
		Round:        input.Round,
		Venue:        input.Venue,
		MatchDate:    matchDate,
		Status:       "scheduled",
	}

	if err := config.DB.Create(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pertandingan"})
		return
	}

	// Load ulang dengan semua relasi
	config.DB.Preload("Tournament").Preload("HomeTeam").Preload("AwayTeam").First(&match, match.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pertandingan berhasil dijadwalkan!",
		"data":    match,
	})
}

// UpdateScore mengupdate skor dan/atau status pertandingan. Hanya admin.
// Endpoint: PATCH /api/matches/:id/score
// Body JSON: { "home_score": 2, "away_score": 1, "status": "finished" }
//
// Dipisahkan dari UpdateMatch agar concern jelas:
//   - UpdateScore  → untuk input hasil pertandingan
//   - UpdateMatch  → untuk ubah jadwal/venue/round
func UpdateScore(c *gin.Context) {
	id := c.Param("id")

	var match models.Match
	if err := config.DB.First(&match, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertandingan tidak ditemukan!"})
		return
	}

	var input struct {
		HomeScore *int   `json:"home_score"`
		AwayScore *int   `json:"away_score"`
		Status    string `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak sesuai!"})
		return
	}

	// Validasi nilai status
	validStatuses := map[string]bool{
		"scheduled": true, "ongoing": true,
		"finished": true, "cancelled": true,
	}
	if input.Status != "" && !validStatuses[input.Status] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":            "Status tidak valid!",
			"allowed_statuses": []string{"scheduled", "ongoing", "finished", "cancelled"},
		})
		return
	}

	// Bangun map update — hanya field yang dikirim
	updates := map[string]interface{}{}
	if input.HomeScore != nil {
		updates["home_score"] = input.HomeScore
	}
	if input.AwayScore != nil {
		updates["away_score"] = input.AwayScore
	}
	if input.Status != "" {
		updates["status"] = input.Status
	}

	config.DB.Model(&match).Updates(updates)

	// Load ulang dengan relasi
	config.DB.Preload("Tournament").Preload("HomeTeam").Preload("AwayTeam").First(&match, id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Skor berhasil diupdate!",
		"data":    match,
	})
}

// UpdateMatch mengupdate detail jadwal pertandingan (round, venue, tanggal). Hanya admin.
// Endpoint: PUT /api/matches/:id
func UpdateMatch(c *gin.Context) {
	id := c.Param("id")

	var match models.Match
	if err := config.DB.First(&match, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertandingan tidak ditemukan!"})
		return
	}

	var input struct {
		Round     string `json:"round"`
		Venue     string `json:"venue"`
		MatchDate string `json:"match_date"`
		Status    string `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak sesuai!"})
		return
	}

	updates := map[string]interface{}{}
	if input.Round != "" {
		updates["round"] = input.Round
	}
	if input.Venue != "" {
		updates["venue"] = input.Venue
	}
	if input.Status != "" {
		updates["status"] = input.Status
	}
	if input.MatchDate != "" {
		matchDate, err := parseDateString(input.MatchDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid!"})
			return
		}
		updates["match_date"] = matchDate
	}

	config.DB.Model(&match).Updates(updates)
	config.DB.Preload("Tournament").Preload("HomeTeam").Preload("AwayTeam").First(&match, id)

	c.JSON(http.StatusOK, gin.H{"message": "Data pertandingan berhasil diupdate!", "data": match})
}

// DeleteMatch menghapus pertandingan. Hanya admin.
// Endpoint: DELETE /api/matches/:id
func DeleteMatch(c *gin.Context) {
	id := c.Param("id")

	var match models.Match
	if err := config.DB.First(&match, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertandingan tidak ditemukan!"})
		return
	}

	config.DB.Delete(&match)
	c.JSON(http.StatusOK, gin.H{"message": "Pertandingan berhasil dihapus!"})
}
