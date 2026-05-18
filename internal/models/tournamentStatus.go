package models

// TournamentStatus adalah tabel lookup/referensi untuk status turnamen.
// Contoh data: {1, "upcoming"}, {2, "ongoing"}, {3, "finished"}, {4, "cancelled"}
type TournamentStatus struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(50);unique;not null" json:"name"`
	// Label adalah nama tampilan yang lebih "cantik" untuk Frontend
	Label string `gorm:"type:varchar(100);not null" json:"label"`
}
