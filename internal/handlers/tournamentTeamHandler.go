package handlers

import (
	"github.com/gofiber/fiber/v2"

	"laga-liga-backend/internal/config"
	"laga-liga-backend/internal/models"
)

func GetTournamentTeams(c *fiber.Ctx) error {
	tournamentID := c.Params("id")

	var tournament models.Tournament
	if err := config.DB.First(&tournament, tournamentID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Turnamen tidak ditemukan!"})
	}

	config.DB.Preload("Teams.Players").First(&tournament, tournamentID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"tournament": fiber.Map{
			"id":   tournament.ID,
			"name": tournament.Name,
		},
		"total_teams": len(tournament.Teams),
		"teams":       tournament.Teams,
	})
}

func RegisterTeamToTournament(c *fiber.Ctx) error {
	tournamentID := c.Params("id")

	var tournament models.Tournament
	if err := config.DB.Preload("Teams").First(&tournament, tournamentID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Turnamen tidak ditemukan!"})
	}

	var input struct {
		TeamID uint `json:"team_id"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Field team_id wajib diisi!"})
	}

	var team models.Team
	if err := config.DB.First(&team, input.TeamID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tim tidak ditemukan!"})
	}

	for _, t := range tournament.Teams {
		if t.ID == team.ID {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Tim sudah terdaftar di turnamen ini!",
			})
		}
	}

	if len(tournament.Teams) >= tournament.MaxTeams {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":     "Turnamen sudah penuh! Kapasitas maksimal tercapai.",
			"max_teams": tournament.MaxTeams,
		})
	}

	if err := config.DB.Model(&tournament).Association("Teams").Append(&team); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mendaftarkan tim ke turnamen"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Tim berhasil didaftarkan ke turnamen!",
		"tournament": fiber.Map{
			"id":   tournament.ID,
			"name": tournament.Name,
		},
		"team": fiber.Map{
			"id":   team.ID,
			"name": team.Name,
			"city": team.City,
		},
		"total_teams_terdaftar": len(tournament.Teams) + 1,
	})
}

func RemoveTeamFromTournament(c *fiber.Ctx) error {
	tournamentID := c.Params("id")
	teamID := c.Params("team_id")

	var tournament models.Tournament
	if err := config.DB.First(&tournament, tournamentID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Turnamen tidak ditemukan!"})
	}

	var team models.Team
	if err := config.DB.First(&team, teamID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tim tidak ditemukan!"})
	}

	if err := config.DB.Model(&tournament).Association("Teams").Delete(&team); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengeluarkan tim dari turnamen"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Tim berhasil dikeluarkan dari turnamen!",
		"tournament": fiber.Map{
			"id":   tournament.ID,
			"name": tournament.Name,
		},
		"team": fiber.Map{
			"id":   team.ID,
			"name": team.Name,
		},
	})
}
