package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"laga-liga-backend/config"
	"laga-liga-backend/models"
)

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
