package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"laga-liga-backend/config"
	"laga-liga-backend/models"
)

// RequireAuth adalah middleware yang memastikan request memiliki JWT Token yang valid.
// Cara kerja:
//  1. Ambil header "Authorization" dari request
//  2. Parse token dan verifikasi tanda tangannya
//  3. Cek apakah token sudah kadaluarsa
//  4. Ambil data user dari DB berdasarkan ID di dalam token
//  5. Simpan data user ke dalam context Gin agar bisa dipakai di controller
func RequireAuth(c *gin.Context) {
	// ── LANGKAH 1: Ambil header Authorization ──────────────────────────────
	// Format header yang benar: "Authorization: Bearer <token>"
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Token tidak ditemukan, silakan login terlebih dahulu",
		})
		return
	}

	// ── LANGKAH 2: Pisahkan kata "Bearer" dari token-nya ──────────────────
	// strings.SplitN memotong string menjadi 2 bagian berdasarkan spasi
	// contoh: "Bearer abc123" → ["Bearer", "abc123"]
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Format token salah! Gunakan format: Bearer <token>",
		})
		return
	}
	tokenString := parts[1]

	// ── LANGKAH 3: Parse dan Verifikasi Token ─────────────────────────────
	// jwt.ParseWithClaims akan:
	//   - Decode token
	//   - Verifikasi tanda tangan dengan secret key
	//   - Cek apakah token sudah expired (karena kita pakai jwt.MapClaims)
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Pastikan algoritma yang digunakan adalah HS256 (bukan algoritma lain)
		// Ini penting untuk mencegah serangan "algorithm confusion"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algoritma tidak valid: %v", token.Header["alg"])
		}
		// Kembalikan secret key untuk memverifikasi tanda tangan token
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	// Jika token tidak valid (expired, tanda tangan salah, dll)
	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Token tidak valid atau sudah kadaluarsa, silakan login ulang",
		})
		return
	}

	// ── LANGKAH 4: Ambil data dari dalam Token (Claims) ───────────────────
	// Claims adalah "isi" dari token, berisi data yang kita simpan saat login
	// Ingat di authController.go kita menyimpan: "sub" = user.ID
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak bisa dibaca"})
		return
	}

	// Ambil user ID dari claim "sub"
	// Catatan: JSON number secara default di-parse sebagai float64
	userIDFloat, ok := (*claims)["sub"].(float64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak berisi ID yang valid"})
		return
	}
	userID := uint(userIDFloat)

	// ── LANGKAH 5: Ambil data User dari Database ──────────────────────────
	// Kita cari user berdasarkan ID yang ada di token
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "User tidak ditemukan",
		})
		return
	}

	// ── LANGKAH 6: Simpan User ke Context Gin ─────────────────────────────
	// Ini seperti "menitipkan" data user ke request,
	// sehingga controller yang datang setelah middleware ini
	// bisa mengambil data user tanpa perlu query DB lagi
	c.Set("currentUser", user)

	// Lanjutkan ke handler/controller berikutnya
	c.Next()
}
