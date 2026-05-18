package handlers

import (
	"github.com/gofiber/fiber/v2"

	"laga-liga-backend/internal/config"
	"laga-liga-backend/internal/models"
)

func GetMatches(c *fiber.Ctx) error {
	var matches []models.Match

	query := config.DB.
		Preload("Tournament").
		Preload("HomeTeam").
		Preload("AwayTeam")

	if tid := c.Query("tournament_id"); tid != "" {
		query = query.Where("tournament_id = ?", tid)
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Order("match_date ASC").Find(&matches)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total": len(matches),
		"data":  matches,
	})
}

func GetMatchByID(c *fiber.Ctx) error {
	var match models.Match
	id := c.Params("id")

	if err := config.DB.
		Preload("Tournament").
		Preload("HomeTeam").
		Preload("AwayTeam").
		First(&match, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pertandingan tidak ditemukan!"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": match})
}

func GetMatchesByTournament(c *fiber.Ctx) error {
	tournamentID := c.Params("id")

	var tournament models.Tournament
	if err := config.DB.First(&tournament, tournamentID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Turnamen tidak ditemukan!"})
	}

	var matches []models.Match
	config.DB.
		Preload("HomeTeam").
		Preload("AwayTeam").
		Where("tournament_id = ?", tournamentID).
		Order("match_date ASC").
		Find(&matches)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"tournament": fiber.Map{"id": tournament.ID, "name": tournament.Name},
		"total":      len(matches),
		"data":       matches,
	})
}

func CreateMatch(c *fiber.Ctx) error {
	var input struct {
		TournamentID uint   `json:"tournament_id"`
		HomeTeamID   uint   `json:"home_team_id"`
		AwayTeamID   uint   `json:"away_team_id"`
		Round        string `json:"round"`
		Venue        string `json:"venue"`
		MatchDate    string `json:"match_date"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Data tidak lengkap: " + err.Error()})
	}

	if input.HomeTeamID == input.AwayTeamID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tim kandang dan tandang tidak boleh sama!"})
	}

	var tournament models.Tournament
	if err := config.DB.First(&tournament, input.TournamentID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Turnamen tidak ditemukan!"})
	}

	var homeTeam, awayTeam models.Team
	if err := config.DB.First(&homeTeam, input.HomeTeamID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tim kandang tidak ditemukan!"})
	}
	if err := config.DB.First(&awayTeam, input.AwayTeamID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tim tandang tidak ditemukan!"})
	}

	matchDate, err := parseDateString(input.MatchDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format tanggal tidak valid! Gunakan: 2026-05-10T15:00:00Z",
		})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat pertandingan"})
	}

	config.DB.Preload("Tournament").Preload("HomeTeam").Preload("AwayTeam").First(&match, match.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pertandingan berhasil dijadwalkan!",
		"data":    match,
	})
}

func UpdateScore(c *fiber.Ctx) error {
	id := c.Params("id")

	var match models.Match
	if err := config.DB.First(&match, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pertandingan tidak ditemukan!"})
	}

	var input struct {
		HomeScore *int   `json:"home_score"`
		AwayScore *int   `json:"away_score"`
		Status    string `json:"status"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format data tidak sesuai!"})
	}

	validStatuses := map[string]bool{
		"scheduled": true, "ongoing": true,
		"finished": true, "cancelled": true,
	}
	if input.Status != "" && !validStatuses[input.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":            "Status tidak valid!",
			"allowed_statuses": []string{"scheduled", "ongoing", "finished", "cancelled"},
		})
	}

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
	config.DB.Preload("Tournament").Preload("HomeTeam").Preload("AwayTeam").First(&match, id)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Skor berhasil diupdate!",
		"data":    match,
	})
}

func UpdateMatch(c *fiber.Ctx) error {
	id := c.Params("id")

	var match models.Match
	if err := config.DB.First(&match, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pertandingan tidak ditemukan!"})
	}

	var input struct {
		Round     string `json:"round"`
		Venue     string `json:"venue"`
		MatchDate string `json:"match_date"`
		Status    string `json:"status"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format data tidak sesuai!"})
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
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format tanggal tidak valid!"})
		}
		updates["match_date"] = matchDate
	}

	config.DB.Model(&match).Updates(updates)
	config.DB.Preload("Tournament").Preload("HomeTeam").Preload("AwayTeam").First(&match, id)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Data pertandingan berhasil diupdate!", "data": match})
}

func DeleteMatch(c *fiber.Ctx) error {
	id := c.Params("id")

	var match models.Match
	if err := config.DB.First(&match, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pertandingan tidak ditemukan!"})
	}

	config.DB.Delete(&match)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Pertandingan berhasil dihapus!"})
}
