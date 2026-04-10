package models

import "time"

// Team merepresentasikan sebuah tim/klub yang bisa ikut turnamen.
type Team struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(150);not null;unique" json:"name"`
	City        string    `gorm:"type:varchar(100)" json:"city"`
	LogoURL     string    `gorm:"type:varchar(255)" json:"logo_url"`
	ManagerName string    `gorm:"type:varchar(100)" json:"manager_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relasi One-to-Many: 1 Tim memiliki banyak Pemain
	// GORM menggunakan nama field ini untuk Preload
	Players []Player `gorm:"foreignKey:TeamID" json:"players,omitempty"`

	// Relasi Many-to-Many: Tim bisa ikut banyak Turnamen
	// GORM akan otomatis buat tabel pivot `tournament_teams`
	Tournaments []Tournament `gorm:"many2many:tournament_teams;" json:"tournaments,omitempty"`
}
