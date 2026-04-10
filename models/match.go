package models

import "time"

// Match merepresentasikan satu pertandingan antara HomeTeam vs AwayTeam
// dalam sebuah Turnamen. Skor diisi setelah pertandingan selesai.
type Match struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TournamentID uint      `gorm:"not null" json:"tournament_id"`
	HomeTeamID   uint      `gorm:"not null" json:"home_team_id"`
	AwayTeamID   uint      `gorm:"not null" json:"away_team_id"`

	// Skor: pointer ke int agar bisa nil (sebelum pertandingan dimulai)
	// Kalau pakai int biasa, nilai 0 ambigu — bisa "belum diisi" atau "skor 0"
	HomeScore *int `gorm:"default:null" json:"home_score"`
	AwayScore *int `gorm:"default:null" json:"away_score"`

	// Status pertandingan: scheduled | ongoing | finished | cancelled
	// Default: scheduled (belum dimulai)
	Status   string    `gorm:"type:varchar(20);default:'scheduled'" json:"status"`

	// Babak/grup: "Grup A", "Grup B", "Perempatfinal", "Semifinal", "Final"
	Round    string    `gorm:"type:varchar(50)" json:"round"`

	// Lokasi dan waktu pertandingan
	Venue    string    `gorm:"type:varchar(200)" json:"venue"`
	MatchDate time.Time `json:"match_date"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// ── Relasi (diisi saat Preload) ────────────────────────────────────────
	// GORM tahu ini FK karena nama field = [NamaField]ID
	Tournament Tournament `gorm:"foreignKey:TournamentID" json:"tournament,omitempty"`
	HomeTeam   Team       `gorm:"foreignKey:HomeTeamID"   json:"home_team,omitempty"`
	AwayTeam   Team       `gorm:"foreignKey:AwayTeamID"   json:"away_team,omitempty"`

	// Relasi One-to-Many: 1 Match memiliki banyak Event (Gol/Kartu)
	Events     []MatchEvent `gorm:"foreignKey:MatchID" json:"events,omitempty"`
}
