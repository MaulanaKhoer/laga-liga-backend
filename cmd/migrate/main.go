package main

import (
	"log"

	"github.com/joho/godotenv"
	"laga-liga-backend/internal/config"
	"laga-liga-backend/internal/models"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: file .env tidak ditemukan, menggunakan nilai default OS.")
	}

	config.ConnectDatabase()

	err = config.DB.AutoMigrate(
		&models.User{},
		&models.TournamentStatus{},
		&models.Sport{},
		&models.Tournament{},
		&models.Team{},
		&models.Player{},
		&models.Match{},
		&models.MatchEvent{},
	)
	if err != nil {
		log.Fatal("Gagal melakukan migrasi database:", err)
	}

	if err := config.DB.Exec(`
		CREATE TABLE IF NOT EXISTS tournament_teams (
			tournament_id BIGINT REFERENCES tournaments(id) ON DELETE CASCADE,
			team_id       BIGINT REFERENCES teams(id) ON DELETE CASCADE,
			PRIMARY KEY (tournament_id, team_id)
		);
	`).Error; err != nil {
		log.Fatal("Gagal membuat tabel tournament_teams:", err)
	}

	log.Println("✅ Migrasi database berhasil!")
}
