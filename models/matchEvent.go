package models

import "time"

// MatchEvent merepresentasikan kejadian penting dalam pertandingan (Gol, Kartu, dll)
type MatchEvent struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MatchID   uint      `gorm:"not null" json:"match_id"`
	PlayerID  uint      `gorm:"not null" json:"player_id"`
	TeamID    uint      `gorm:"not null" json:"team_id"` // Tim dari pemain saat kejadian
	
	// Type: goal | yellow_card | red_card | own_goal
	Type      string    `gorm:"type:varchar(20);not null" json:"type"`
	
	// Minute: Menit terjadinya kejadian (misal: 45)
	Minute    int       `gorm:"not null" json:"minute"`
	
	// Note: Catatan tambahan (misal: "Penalti", "Assist by X")
	Note      string    `gorm:"type:varchar(255)" json:"note"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relasi
	Match  Match  `gorm:"foreignKey:MatchID" json:"match,omitempty"`
	Player Player `gorm:"foreignKey:PlayerID" json:"player,omitempty"`
	Team   Team   `gorm:"foreignKey:TeamID" json:"team,omitempty"`
}
