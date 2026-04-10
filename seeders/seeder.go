package seeders

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"laga-liga-backend/models"
)

// SeedAll menjalankan semua seeder secara berurutan.
// Aman dijalankan berulang kali — tidak akan duplikat data.
func SeedAll(db *gorm.DB) {
	log.Println("🌱 Mulai proses seeding database...")
	seedTournamentStatuses(db) // ← Harus pertama! Tournament butuh status_id
	seedUsers(db)
	seedTournaments(db)
	log.Println("✅ Seeding selesai!")
}

// seedTournamentStatuses mengisi tabel lookup status turnamen
func seedTournamentStatuses(db *gorm.DB) {
	statuses := []models.TournamentStatus{
		{ID: 1, Name: "upcoming", Label: "Akan Datang"},
		{ID: 2, Name: "ongoing", Label: "Sedang Berlangsung"},
		{ID: 3, Name: "finished", Label: "Selesai"},
		{ID: 4, Name: "cancelled", Label: "Dibatalkan"},
	}

	for _, s := range statuses {
		var count int64
		db.Model(&models.TournamentStatus{}).Where("id = ?", s.ID).Count(&count)
		if count > 0 {
			fmt.Printf("   ⏭️  Status '%s' sudah ada, skip.\n", s.Name)
			continue
		}
		if err := db.Create(&s).Error; err != nil {
			log.Printf("   ❌ Gagal membuat status '%s': %v\n", s.Name, err)
		} else {
			fmt.Printf("   ✅ Status [%d] '%s' (%s) berhasil dibuat.\n", s.ID, s.Name, s.Label)
		}
	}
}

// seedUsers membuat akun dummy: 1 admin dan 2 user biasa
func seedUsers(db *gorm.DB) {
	users := []struct {
		Name     string
		Email    string
		Password string
		Role     string
	}{
		{"Admin Laga-Liga", "admin@lagaliga.id", "admin123", "admin"},
		{"Budi Santoso", "budi@email.com", "password123", "user"},
		{"Sari Dewi", "sari@email.com", "password123", "user"},
	}

	for _, u := range users {
		var count int64
		db.Model(&models.User{}).Where("email = ?", u.Email).Count(&count)
		if count > 0 {
			fmt.Printf("   ⏭️  User '%s' sudah ada, skip.\n", u.Email)
			continue
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("   ❌ Gagal hash password untuk %s: %v\n", u.Email, err)
			continue
		}

		user := models.User{
			Name:     u.Name,
			Email:    u.Email,
			Password: string(hashed),
			Role:     u.Role,
		}

		if err := db.Create(&user).Error; err != nil {
			log.Printf("   ❌ Gagal membuat user %s: %v\n", u.Email, err)
		} else {
			fmt.Printf("   ✅ User '%s' (%s) berhasil dibuat.\n", u.Name, u.Role)
		}
	}
}

// seedTournaments membuat data turnamen dummy
// StatusID mengacu ke tabel tournament_statuses
func seedTournaments(db *gorm.DB) {
	tournaments := []models.Tournament{
		{
			Name:        "Piala Indonesia 2026",
			Description: "Turnamen sepak bola nasional antar klub Indonesia",
			Location:    "Jakarta, Indonesia",
			StartDate:   time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2026, 8, 31, 0, 0, 0, 0, time.UTC),
			MaxTeams:    16,
			StatusID:    1, // upcoming
		},
		{
			Name:        "Liga Futsal Bandung Open",
			Description: "Kompetisi futsal terbuka untuk umum se-Jawa Barat",
			Location:    "Bandung, Jawa Barat",
			StartDate:   time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC),
			MaxTeams:    8,
			StatusID:    2, // ongoing
		},
		{
			Name:        "Turnamen Mini Soccer Surabaya Cup",
			Description: "Turnamen mini soccer tahunan untuk kategori usia 17-25 tahun",
			Location:    "Surabaya, Jawa Timur",
			StartDate:   time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2026, 6, 20, 0, 0, 0, 0, time.UTC),
			MaxTeams:    12,
			StatusID:    1, // upcoming
		},
	}

	for _, t := range tournaments {
		var count int64
		db.Model(&models.Tournament{}).Where("name = ?", t.Name).Count(&count)
		if count > 0 {
			fmt.Printf("   ⏭️  Turnamen '%s' sudah ada, skip.\n", t.Name)
			continue
		}
		if err := db.Create(&t).Error; err != nil {
			log.Printf("   ❌ Gagal membuat turnamen '%s': %v\n", t.Name, err)
		} else {
			fmt.Printf("   ✅ Turnamen '%s' (status_id: %d) berhasil dibuat.\n", t.Name, t.StatusID)
		}
	}
}
