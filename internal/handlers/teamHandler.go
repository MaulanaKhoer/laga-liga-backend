package handlers

import (
	"github.com/gofiber/fiber/v2"

	"laga-liga-backend/internal/config"
	"laga-liga-backend/internal/models"
)

func GetTeams(c *fiber.Ctx) error {
	var teams []models.Team
	config.DB.Preload("Players").Find(&teams)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": teams})
}

func GetTeamByID(c *fiber.Ctx) error {
	var team models.Team
	id := c.Params("id")

	if err := config.DB.Preload("Players").First(&team, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tim tidak ditemukan!"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": team})
}

func CreateTeam(c *fiber.Ctx) error {
	var input models.Team
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format data tidak sesuai!"})
	}

	if err := config.DB.Create(&input).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan tim"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Tim berhasil dibuat!", "data": input})
}

func UpdateTeam(c *fiber.Ctx) error {
	var team models.Team
	id := c.Params("id")

	if err := config.DB.First(&team, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tim tidak ditemukan!"})
	}

	var input models.Team
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format data tidak sesuai!"})
	}

	config.DB.Model(&team).Updates(input)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Tim berhasil diupdate!", "data": team})
}

func DeleteTeam(c *fiber.Ctx) error {
	var team models.Team
	id := c.Params("id")

	if err := config.DB.First(&team, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tim tidak ditemukan!"})
	}

	config.DB.Delete(&team)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Tim berhasil dihapus!"})
}

func AddPlayerToTeam(c *fiber.Ctx) error {
	teamID := c.Params("id")

	var team models.Team
	if err := config.DB.First(&team, teamID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tim tidak ditemukan!"})
	}

	var input models.Player
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format data tidak sesuai!"})
	}

	input.TeamID = team.ID

	if err := config.DB.Create(&input).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menambahkan pemain"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pemain berhasil ditambahkan ke tim!",
		"data":    input,
	})
}
