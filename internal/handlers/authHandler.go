package handlers

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"laga-liga-backend/internal/config"
	"laga-liga-backend/internal/models"
)

type AuthInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
	var input AuthInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email dan Password wajib ada!"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengenkripsi kata sandi"})
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email mungkin sudah digunakan!"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Registrasi berhasil, silakan login!"})
}

func Login(c *fiber.Ctx) error {
	var input AuthInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format salah!"})
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email atau Password salah!"})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email atau Password salah!"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat token login"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login sukses",
		"token":   tokenString,
	})
}

func GetMe(c *fiber.Ctx) error {
	userRaw := c.Locals("currentUser")
	if userRaw == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Tidak bisa mengambil data user"})
	}

	user, ok := userRaw.(models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Format data user tidak valid"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data profil berhasil diambil",
		"user":    user,
	})
}
