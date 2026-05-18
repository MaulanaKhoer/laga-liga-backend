package models

import "time"

// Sport adalah tabel referensi untuk cabang olahraga atau aplikasi game.
type Sport struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null;unique" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	IconPath    string    `gorm:"type:varchar(255)" json:"icon_path"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
