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
	seedTeams(db)
	seedPlayers(db)
	seedTournamentTeams(db) // ← Daftarkan tim ke turnamen
	seedMatches(db)         // ← Berikan jadwal pertandingan
	seedMatchEvents(db)    // ← Detail gol dan kartu
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

// seedTeams membuat data tim dummy
func seedTeams(db *gorm.DB) {
	teams := []models.Team{
		{Name: "Persija Jakarta", City: "Jakarta", ManagerName: "Budi Gunawan"},
		{Name: "Persib Bandung", City: "Bandung", ManagerName: "Agus Setiawan"},
		{Name: "Arema FC", City: "Malang", ManagerName: "Slamet Santoso"},
		{Name: "PSM Makassar", City: "Makassar", ManagerName: "Rahman Dali"},
	}

	for _, t := range teams {
		var count int64
		db.Model(&models.Team{}).Where("name = ?", t.Name).Count(&count)
		if count > 0 {
			fmt.Printf("   ⏭️  Tim '%s' sudah ada, skip.\n", t.Name)
			continue
		}
		if err := db.Create(&t).Error; err != nil {
			log.Printf("   ❌ Gagal membuat tim '%s': %v\n", t.Name, err)
		} else {
			fmt.Printf("   ✅ Tim '%s' berhasil dibuat.\n", t.Name)
		}
	}
}

// seedPlayers membuat pemain dummy untuk tiap tim
func seedPlayers(db *gorm.DB) {
	// Ambil ID tim yang sudah ada
	var teams []models.Team
	db.Find(&teams)
	if len(teams) == 0 {
		log.Println("   ⚠️  Tidak ada tim, skip seed pemain.")
		return
	}

	// Untuk setiap tim, buat 3 pemain dummy
	positions := []string{"GK", "CB", "ST"}
	namePrefix := [][]string{
		{"Andi", "Budi", "Candra"},
		{"Deni", "Eko", "Fajar"},
		{"Galih", "Hendra", "Ivan"},
		{"Joko", "Kevin", "Luthfi"},
	}

	for i, team := range teams {
		var count int64
		db.Model(&models.Player{}).Where("team_id = ?", team.ID).Count(&count)
		if count > 0 {
			fmt.Printf("   ⏭️  Pemain tim '%s' sudah ada, skip.\n", team.Name)
			continue
		}

		for j, pos := range positions {
			player := models.Player{
				Name:         namePrefix[i][j] + " " + team.City,
				Position:     pos,
				JerseyNumber: j + 1,
				TeamID:       team.ID,
			}
			db.Create(&player)
		}
		fmt.Printf("   ✅ 3 pemain untuk tim '%s' berhasil dibuat.\n", team.Name)
	}
}

// seedTournamentTeams mendaftarkan tim yang ada ke turnamen dummy
func seedTournamentTeams(db *gorm.DB) {
	var tournament models.Tournament
	if err := db.Where("name = ?", "Piala Indonesia 2026").First(&tournament).Error; err != nil {
		return
	}

	var teams []models.Team
	db.Find(&teams)

	// Cek apakah sudah ada tim terdaftar
	var count int64
	db.Table("tournament_teams").Where("tournament_id = ?", tournament.ID).Count(&count)
	if count > 0 {
		fmt.Printf("   ⏭️  Tim sudah terdaftar di '%s', skip.\n", tournament.Name)
		return
	}

	// Daftarkan semua tim hasil seeder ke turnamen ini
	for _, team := range teams {
		db.Model(&tournament).Association("Teams").Append(&team)
	}
	fmt.Printf("   ✅ %d tim berhasil didaftarkan ke '%s'.\n", len(teams), tournament.Name)
}

// seedMatches membuat data pertandingan dummy
func seedMatches(db *gorm.DB) {
	var tournament models.Tournament
	if err := db.Where("name = ?", "Piala Indonesia 2026").First(&tournament).Error; err != nil {
		return
	}

	// Ambil 2 tim pertama untuk buat pertandingan contoh
	var teams []models.Team
	db.Limit(2).Find(&teams)
	if len(teams) < 2 {
		return
	}

	score1, score2 := 2, 1
	matches := []models.Match{
		{
			TournamentID: tournament.ID,
			HomeTeamID:   teams[0].ID,
			AwayTeamID:   teams[1].ID,
			HomeScore:    &score1,
			AwayScore:    &score2,
			Status:       "finished",
			Round:        "Grup A",
			Venue:        "Stadion Utama GBK",
			MatchDate:    time.Now().AddDate(0, 0, -1), // Kemarin
		},
		{
			TournamentID: tournament.ID,
			HomeTeamID:   teams[1].ID,
			AwayTeamID:   teams[0].ID,
			Status:       "scheduled",
			Round:        "Grup A",
			Venue:        "Stadion Si Jalak Harupat",
			MatchDate:    time.Now().AddDate(0, 0, 7), // Minggu depan
		},
	}

	for _, m := range matches {
		var count int64
		db.Model(&models.Match{}).Where(
			"tournament_id = ? AND home_team_id = ? AND away_team_id = ? AND round = ?",
			m.TournamentID, m.HomeTeamID, m.AwayTeamID, m.Round,
		).Count(&count)

		if count > 0 {
			continue
		}

		db.Create(&m)
	}
	fmt.Println("   ✅ Data pertandingan dummy berhasil dibuat.")
}

// seedMatchEvents mengisi detail gol untuk pertandingan yang sudah selesai
func seedMatchEvents(db *gorm.DB) {
	// Ambil pertandingan yang sudah selesai
	var matches []models.Match
	if err := db.Where("status = ?", "finished").Find(&matches).Error; err != nil {
		return
	}

	for _, m := range matches {
		// Cek apakah sudah ada event
		var count int64
		db.Model(&models.MatchEvent{}).Where("match_id = ?", m.ID).Count(&count)
		if count > 0 {
			continue
		}

		// Ambil pemain dari kedua tim
		var homePlayers, awayPlayers []models.Player
		db.Where("team_id = ?", m.HomeTeamID).Find(&homePlayers)
		db.Where("team_id = ?", m.AwayTeamID).Find(&awayPlayers)

		if len(homePlayers) == 0 || len(awayPlayers) == 0 {
			continue
		}

		// Simulasi gol sesuai skor match (misal 2-1)
		hScore := 0
		if m.HomeScore != nil {
			hScore = *m.HomeScore
		}
		aScore := 0
		if m.AwayScore != nil {
			aScore = *m.AwayScore
		}

		// Home goals
		for i := 0; i < hScore; i++ {
			event := models.MatchEvent{
				MatchID:  m.ID,
				PlayerID: homePlayers[0].ID, // Anggap striker yang cetak gol
				TeamID:   m.HomeTeamID,
				Type:     "goal",
				Minute:   10 + (i * 20),
			}
			db.Create(&event)
		}

		// Away goals
		for i := 0; i < aScore; i++ {
			event := models.MatchEvent{
				MatchID:  m.ID,
				PlayerID: awayPlayers[0].ID,
				TeamID:   m.AwayTeamID,
				Type:     "goal",
				Minute:   35,
			}
			db.Create(&event)
		}
	}
	fmt.Println("   ✅ Detail kejadian pertandingan berhasil dibuat.")
}
