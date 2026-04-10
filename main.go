package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"laga-liga-backend/config"
	"laga-liga-backend/models"
	"laga-liga-backend/routes"
	"laga-liga-backend/seeders"
)

func main() {
	// 1. Load file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: file .env tidak ditemukan, menggunakan nilai default OS.")
	}

	// 2. Koneksi ke Database
	config.ConnectDatabase()

	// 3. AutoMigrate: GORM sinkronisasi struct ke tabel PostgreSQL
	// URUTAN PENTING: TournamentStatus harus dibuat SEBELUM Tournament
	// karena Tournament punya foreign key ke tournament_statuses
	config.DB.Migrator().DropTable(&models.Tournament{})
	err = config.DB.AutoMigrate(
		&models.User{},
		&models.TournamentStatus{}, // ← tabel referensi dulu
		&models.Tournament{},       // ← baru tabel yang punya FK
	)
	if err != nil {
		log.Fatal("Gagal melakukan migrasi database:", err)
	}

	// 4. Jalankan Seeder (isi data dummy jika belum ada)
	seeders.SeedAll(config.DB)

	// 5. Setup Router Gin
	r := gin.Default()

	// 6. Daftarkan semua routes (dipindah ke routes/routes.go agar main.go ringkas)
	routes.SetupRoutes(r)

	// 7. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("🚀 Server berjalan di http://localhost:%s", port)
	r.Run(":" + port)
}
