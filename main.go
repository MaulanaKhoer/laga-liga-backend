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
	"laga-liga-backend/middleware"
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
	// URUTAN PENTING:
	//   - TournamentStatus harus sebelum Tournament (FK status_id)
	//   - Team harus sebelum Player (FK team_id)
	//   - Tournament & Team harus ada sebelum pivot tournament_teams dibuat
	err = config.DB.AutoMigrate(
		&models.User{},
		&models.TournamentStatus{},
		&models.Tournament{},
		&models.Team{},   // ← Team sebelum Player (Player punya FK ke Team)
		&models.Player{}, // ← Player punya FK ke Team
		&models.Match{},  // ← Tabel Match baru!
		&models.MatchEvent{}, // ← Tabel Detail Pertandingan
	)
	if err != nil {
		log.Fatal("Gagal melakukan migrasi database:", err)
	}

	// Migrate tabel pivot many2many: tournament_teams
	// Tabel ini dibuat otomatis GORM dari tag `gorm:"many2many:tournament_teams;"`
	// di model Team. Kita perlu pastikan tabel ini ada terpisah dari AutoMigrate.
	if err := config.DB.Exec(`
		CREATE TABLE IF NOT EXISTS tournament_teams (
			tournament_id BIGINT REFERENCES tournaments(id) ON DELETE CASCADE,
			team_id       BIGINT REFERENCES teams(id) ON DELETE CASCADE,
			PRIMARY KEY (tournament_id, team_id)
		);
	`).Error; err != nil {
		log.Fatal("Gagal membuat tabel tournament_teams:", err)
	}

	// 4. Jalankan Seeder (isi data dummy jika belum ada)
	seeders.SeedAll(config.DB)

	// 5. Setup Router Gin
	r := gin.Default()

	// Gunakan middleware CORS
	r.Use(middleware.SetupCORS())

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
