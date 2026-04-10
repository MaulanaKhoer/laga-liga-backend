package models

import "time"

// Tournament adalah tabel utama untuk data turnamen.
// StatusID adalah Foreign Key yang merujuk ke tabel tournament_statuses.
type Tournament struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(150);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Location    string    `gorm:"type:varchar(200)" json:"location"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	MaxTeams    int       `gorm:"default:8" json:"max_teams"`

	// Foreign Key ke tabel tournament_statuses
	// `gorm:"default:1"` artinya status default adalah ID=1 (upcoming)
	StatusID uint `gorm:"not null;default:1" json:"status_id"`

	// Relasi: GORM akan otomatis JOIN ke tabel tournament_statuses
	// Field ini diisi oleh GORM saat kita pakai Preload("Status")
	// `json:"status"` agar Frontend dapat object status lengkap, bukan hanya ID-nya
	Status TournamentStatus `gorm:"foreignKey:StatusID" json:"status"`

	// Relasi Many-to-Many: Turnamen bisa diikuti banyak Tim
	// GORM menggunakan tabel pivot `tournament_teams` (sama dengan di Team)
	Teams []Team `gorm:"many2many:tournament_teams;" json:"teams,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
