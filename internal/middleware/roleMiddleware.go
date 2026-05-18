package middleware

import (
	"github.com/gofiber/fiber/v2"

	"laga-liga-backend/internal/models"
)

func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRaw := c.Locals("currentUser")
		if userRaw == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Tidak terautentikasi",
			})
		}

		user, ok := userRaw.(models.User)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal membaca data user",
			})
		}

		for _, role := range allowedRoles {
			if user.Role == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Akses ditolak! Kamu tidak memiliki izin untuk melakukan aksi ini.",
		})
	}
}
