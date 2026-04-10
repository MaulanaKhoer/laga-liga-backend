package controllers

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	"laga-liga-backend/config"
	"laga-liga-backend/models"
)

type StandingRow struct {
	Rank      int    `json:"rank"`
	TeamID    uint   `json:"team_id"`
	TeamName  string `json:"team_name"`
	Played    int    `json:"played"`
	Won       int    `json:"won"`
	Draw      int    `json:"draw"`
	Lost      int    `json:"lost"`
	GF        int    `json:"gf"` // Goals For (memasang gol)
	GA        int    `json:"ga"` // Goals Against (kebobolan)
	GD        int    `json:"gd"` // Goal Difference (selisih gol)
	Points    int    `json:"points"`
}

// ScorerRow merepresentasikan data pencetak gol
type ScorerRow struct {
	PlayerID   uint   `json:"player_id"`
	PlayerName string `json:"player_name"`
	TeamName   string `json:"team_name"`
	Goals      int    `json:"goals"`
}

// CreateTournament membuat turnamen baru. Hanya admin.
func CreateTournament(c *gin.Context) {
	var input models.Tournament
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak sesuai!"})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan ke database"})
		return
	}

	// Setelah create, load ulang dengan relasi Status agar response lengkap
	config.DB.Preload("Status").First(&input, input.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Turnamen berhasil dibuat!",
		"data":    input,
	})
}

// GetTournaments mengambil semua turnamen beserta data statusnya.
func GetTournaments(c *gin.Context) {
	var tournaments []models.Tournament

	// Preload("Status") artinya: untuk setiap tournament,
	// GORM otomatis JOIN ke tabel tournament_statuses dan isi field Status.
	// Tanpa Preload, field Status akan kosong (zero value).
	config.DB.Preload("Status").Find(&tournaments)

	c.JSON(http.StatusOK, gin.H{"data": tournaments})
}

// GetTournamentByID mengambil 1 turnamen berdasarkan ID.
func GetTournamentByID(c *gin.Context) {
	var tournament models.Tournament
	id := c.Param("id")

	if err := config.DB.Preload("Status").First(&tournament, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Turnamen tidak ditemukan!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tournament})
}

// UpdateTournament mengupdate data turnamen. Hanya admin.
func UpdateTournament(c *gin.Context) {
	var tournament models.Tournament
	id := c.Param("id")

	if err := config.DB.First(&tournament, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Turnamen tidak ditemukan!"})
		return
	}

	// Gunakan struct terpisah untuk input agar StatusID bisa di-update
	var input models.Tournament
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak sesuai!"})
		return
	}

	config.DB.Model(&tournament).Updates(input)

	// Load ulang dengan relasi Status
	config.DB.Preload("Status").First(&tournament, id)

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diupdate!", "data": tournament})
}

// DeleteTournament menghapus turnamen. Hanya admin.
func DeleteTournament(c *gin.Context) {
	var tournament models.Tournament
	id := c.Param("id")

	if err := config.DB.First(&tournament, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Turnamen tidak ditemukan!"})
		return
	}

	config.DB.Delete(&tournament)
	c.JSON(http.StatusOK, gin.H{"message": "Turnamen berhasil dihapus!"})
}

// GetTournamentStandings menghitung klasemen turnamen secara dinamis (on-the-fly).
// Endpoint: GET /api/tournaments/:id/standings
func GetTournamentStandings(c *gin.Context) {
	tournamentID := c.Param("id")

	// 1. Ambil data turnamen beserta tim-tim yang terdaftar
	var tournament models.Tournament
	if err := config.DB.Preload("Teams").First(&tournament, tournamentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Turnamen tidak ditemukan!"})
		return
	}

	// 2. Ambil semua pertandingan yang sudah SELESAI di turnamen ini
	var matches []models.Match
	config.DB.Where("tournament_id = ? AND status = ?", tournamentID, "finished").Find(&matches)

	// 3. Siapkan map untuk menampung statistik tiap tim
	standingsMap := make(map[uint]*StandingRow)
	for _, team := range tournament.Teams {
		standingsMap[team.ID] = &StandingRow{
			TeamID:   team.ID,
			TeamName: team.Name,
		}
	}

	// 4. Hitung statistik dari setiap pertandingan
	for _, m := range matches {
		homeRow := standingsMap[m.HomeTeamID]
		awayRow := standingsMap[m.AwayTeamID]

		// Jika ada match untuk tim yang tidak terdaftar di tournament_teams (data anomali), skip
		if homeRow == nil || awayRow == nil {
			continue
		}

		homeRow.Played++
		awayRow.Played++

		hScore := 0
		if m.HomeScore != nil {
			hScore = *m.HomeScore
		}
		aScore := 0
		if m.AwayScore != nil {
			aScore = *m.AwayScore
		}

		homeRow.GF += hScore
		homeRow.GA += aScore

		awayRow.GF += aScore
		awayRow.GA += hScore

		if hScore > aScore {
			// Home Menang
			homeRow.Won++
			homeRow.Points += 3
			awayRow.Lost++
		} else if hScore < aScore {
			// Away Menang
			awayRow.Won++
			awayRow.Points += 3
			homeRow.Lost++
		} else {
			// Seri
			homeRow.Draw++
			awayRow.Draw++
			homeRow.Points += 1
			awayRow.Points += 1
		}
	}

	// 5. Konversi map ke slice agar bisa disortir
	var standings []StandingRow
	for _, row := range standingsMap {
		row.GD = row.GF - row.GA
		standings = append(standings, *row)
	}

	// 6. Sortir klasemen: Poin -> GD -> GF
	sort.Slice(standings, func(i, j int) bool {
		if standings[i].Points != standings[j].Points {
			return standings[i].Points > standings[j].Points
		}
		if standings[i].GD != standings[j].GD {
			return standings[i].GD > standings[j].GD
		}
		return standings[i].GF > standings[j].GF
	})

	// 7. Berikan nomor peringkat (Rank)
	for i := range standings {
		standings[i].Rank = i + 1
	}

	c.JSON(http.StatusOK, gin.H{
		"tournament": gin.H{"id": tournament.ID, "name": tournament.Name},
		"data":       standings,
	})
}

// GetTopScorers mengambil daftar pencetak gol terbanyak dalam satu turnamen.
// Endpoint: GET /api/tournaments/:id/top-scorers
func GetTopScorers(c *gin.Context) {
	tournamentID := c.Param("id")

	// 1. Ambil semua gol dari turnamen ini
	// Caranya: JOIN MatchEvent -> Match, filter by TournamentID
	var scorers []ScorerRow
	
	// Query GORM:
	// SELECT p.id as player_id, p.name as player_name, t.name as team_name, count(me.id) as goals
	// FROM match_events me
	// JOIN players p ON p.id = me.player_id
	// JOIN teams t ON t.id = me.team_id
	// JOIN matches m ON m.id = me.match_id
	// WHERE m.tournament_id = ? AND me.type = 'goal'
	// GROUP BY p.id, p.name, t.name
	// ORDER BY goals DESC
	
	err := config.DB.Table("match_events").
		Select("players.id as player_id, players.name as player_name, teams.name as team_name, count(match_events.id) as goals").
		Joins("JOIN players ON players.id = match_events.player_id").
		Joins("JOIN teams ON teams.id = match_events.team_id").
		Joins("JOIN matches ON matches.id = match_events.match_id").
		Where("matches.tournament_id = ? AND match_events.type = ?", tournamentID, "goal").
		Group("players.id, players.name, teams.name").
		Order("goals DESC").
		Scan(&scorers).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data top scorer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tournament_id": tournamentID,
		"data":          scorers,
	})
}
