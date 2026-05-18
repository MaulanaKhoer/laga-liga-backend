package models

import "time"

// Player merepresentasikan pemain yang tergabung dalam sebuah tim.
type Player struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(150);not null" json:"name"`
	Position    string    `gorm:"type:varchar(50)" json:"position"` // GK, CB, CM, ST, dll
	JerseyNumber int      `gorm:"" json:"jersey_number"`
	DateOfBirth  string   `gorm:"type:varchar(20)" json:"date_of_birth"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Foreign Key: setiap Pemain pasti milik satu Tim
	TeamID uint `gorm:"not null" json:"team_id"`

	// Relasi BelongsTo: GORM isi field ini saat Preload("Team")
	Team Team `gorm:"foreignKey:TeamID" json:"team,omitempty"`
}
