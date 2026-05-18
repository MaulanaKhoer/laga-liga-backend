package handlers

import (
	"github.com/gofiber/fiber/v2"

	"laga-liga-backend/internal/config"
	"laga-liga-backend/internal/models"
)

func GetSports(c *fiber.Ctx) error {
	var sports []models.Sport
	config.DB.Order("name ASC").Find(&sports)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total": len(sports),
		"data":  sports,
	})
}
