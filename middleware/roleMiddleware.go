package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"laga-liga-backend/models"
)

// RequireRole adalah middleware GENERATOR — fungsi yang menghasilkan fungsi middleware.
//
// Kenapa dibuat seperti ini? Karena kita ingin bisa menulis:
//   protected.POST("/tournaments", middleware.RequireRole("admin"), controllers.Create)
//
// Dengan cara ini, kita bisa menggunakan 1 fungsi untuk berbagai role:
//   middleware.RequireRole("admin")
//   middleware.RequireRole("moderator")
//   middleware.RequireRole("admin", "moderator")  ← bisa multiple role!
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	// Fungsi ini mengembalikan fungsi middleware yang sesungguhnya
	return func(c *gin.Context) {
		// ── STEP 1: Ambil data user dari context (diisi oleh RequireAuth) ──
		// PENTING: RequireRole harus selalu dipakai SETELAH RequireAuth,
		// karena RequireAuth-lah yang menyimpan "currentUser" ke context.
		userRaw, exists := c.Get("currentUser")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Tidak terautentikasi",
			})
			return
		}

		// ── STEP 2: Konversi interface{} ke models.User ────────────────────
		user, ok := userRaw.(models.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal membaca data user",
			})
			return
		}

		// ── STEP 3: Cek apakah role user ada di daftar role yang diizinkan ──
		// Kita loop allowedRoles dan cek satu per satu
		for _, role := range allowedRoles {
			if user.Role == role {
				// Role cocok! Lanjutkan ke handler berikutnya
				c.Next()
				return
			}
		}

		// ── STEP 4: Jika tidak ada role yang cocok, tolak request ──────────
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			// 403 Forbidden: user sudah login, tapi tidak punya hak akses
			"error": "Akses ditolak! Kamu tidak memiliki izin untuk melakukan aksi ini.",
		})
	}
}
