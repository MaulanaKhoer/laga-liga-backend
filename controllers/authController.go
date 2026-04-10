package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"laga-liga-backend/config"
	"laga-liga-backend/models"
)

// Struktur bantuan untuk menangkap JSON input (karena model asli tak mewajibkan password untuk di-return)
type AuthInput struct {
	Name     string `json:"name"` 
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input AuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email dan Password wajib ada!"})
		return
	}

	// 1. Hash Password menggunakan bcrypt agar tidak tersimpan dalam format teks biasa
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi kata sandi"})
		return
	}

	// 2. Buat User baru dengan password yang sudah di-hash
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	// 3. Simpan ke Database
	if err := config.DB.Create(&user).Error; err != nil {
		// Biasanya error karena email sudah terdaftar (terkena aturan UNIQUE constraint dari PostgreSQL)
		c.JSON(http.StatusConflict, gin.H{"error": "Email mungkin sudah digunakan!"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registrasi berhasil, silakan login!"})
}

func Login(c *gin.Context) {
	var input AuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format salah!"})
		return
	}

	var user models.User
	// 1. Cari user berdasarkan Email dari DB
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau Password salah!"})
		return
	}

	// 2. Cek apakah password cocok dengan Hash di database
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau Password salah!"})
		return
	}

	// 3. Jika Valid, buatkan JWT Token yang berlaku 24 Jam
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Tandatangani token pakai secret key dari .env
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token login"})
		return
	}

	// 4. Return token ke Frontend (Vue/React) agar mereka bisa menyimpannya di localStorage/cookie
	c.JSON(http.StatusOK, gin.H{
		"message": "Login sukses",
		"token":   tokenString,
	})
}

// GetMe mengembalikan profil user yang sedang login.
// Endpoint ini dilindungi middleware, jadi kita PASTI sudah punya "currentUser" di context.
func GetMe(c *gin.Context) {
	// Ambil data user dari context yang sudah diisi oleh middleware RequireAuth
	// c.MustGet() akan panic jika key tidak ada — tapi karena ada middleware, ini aman
	userRaw, exists := c.Get("currentUser")
	if !exists {
		// Seharusnya tidak pernah terjadi jika middleware terpasang benar
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tidak bisa mengambil data user"})
		return
	}

	// Type assertion: ubah interface{} kembali ke tipe models.User
	// Ini diperlukan karena c.Get() mengembalikan tipe interface{} (tipe umum)
	user, ok := userRaw.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Format data user tidak valid"})
		return
	}

	// Kembalikan data user (field Password otomatis disembunyikan karena json:"-" di model)
	c.JSON(http.StatusOK, gin.H{
		"message": "Data profil berhasil diambil",
		"user":    user,
	})
}
