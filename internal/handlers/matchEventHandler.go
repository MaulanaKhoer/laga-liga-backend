package handlers

import (
	"github.com/gofiber/fiber/v2"

	"laga-liga-backend/internal/config"
	"laga-liga-backend/internal/models"
)

func AddMatchEvent(c *fiber.Ctx) error {
	matchID := c.Params("id")

	var match models.Match
	if err := config.DB.First(&match, matchID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pertandingan tidak ditemukan!"})
	}

	var input struct {
		PlayerID uint   `json:"player_id"`
		Type     string `json:"type"` 
		Minute   int    `json:"minute"`
		Note     string `json:"note"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Data tidak lengkap: " + err.Error()})
	}

	var player models.Player
	if err := config.DB.First(&player, input.PlayerID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pemain tidak ditemukan!"})
	}

	if player.TeamID != match.HomeTeamID && player.TeamID != match.AwayTeamID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Pemain ini tidak terdaftar di tim yang sedang bertanding!",
		})
	}

	event := models.MatchEvent{
		MatchID:  match.ID,
		PlayerID: player.ID,
		TeamID:   player.TeamID,
		Type:     input.Type,
		Minute:   input.Minute,
		Note:     input.Note,
	}

	if err := config.DB.Create(&event).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan event"})
	}

	config.DB.Preload("Player").Preload("Team").First(&event, event.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Event berhasil dicatat!",
		"data":    event,
	})
}

func GetMatchEvents(c *fiber.Ctx) error {
	matchID := c.Params("id")

	var events []models.MatchEvent
	config.DB.
		Preload("Player").
		Preload("Team").
		Where("match_id = ?", matchID).
		Order("minute ASC").
		Find(&events)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": events})
}

func DeleteMatchEvent(c *fiber.Ctx) error {
	eventID := c.Params("event_id")

	var event models.MatchEvent
	if err := config.DB.First(&event, eventID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event tidak ditemukan!"})
	}

	config.DB.Delete(&event)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Event berhasil dihapus!"})
}
