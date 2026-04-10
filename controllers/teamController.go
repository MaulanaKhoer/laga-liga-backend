package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"laga-liga-backend/config"
	"laga-liga-backend/models"
)

// GetTeams mengambil semua tim beserta pemainnya
func GetTeams(c *gin.Context) {
	var teams []models.Team

	// Preload("Players") → untuk setiap tim, ikutkan data pemainnya
	config.DB.Preload("Players").Find(&teams)

	c.JSON(http.StatusOK, gin.H{"data": teams})
}

// GetTeamByID mengambil 1 tim + semua pemainnya
func GetTeamByID(c *gin.Context) {
	var team models.Team
	id := c.Param("id")

	if err := config.DB.Preload("Players").First(&team, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tim tidak ditemukan!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": team})
}

// CreateTeam membuat tim baru. Hanya admin.
func CreateTeam(c *gin.Context) {
	var input models.Team
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak sesuai!"})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan tim"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Tim berhasil dibuat!", "data": input})
}

// UpdateTeam mengupdate data tim. Hanya admin.
func UpdateTeam(c *gin.Context) {
	var team models.Team
	id := c.Param("id")

	if err := config.DB.First(&team, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tim tidak ditemukan!"})
		return
	}

	var input models.Team
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak sesuai!"})
		return
	}

	config.DB.Model(&team).Updates(input)
	c.JSON(http.StatusOK, gin.H{"message": "Tim berhasil diupdate!", "data": team})
}

// DeleteTeam menghapus tim. Hanya admin.
func DeleteTeam(c *gin.Context) {
	var team models.Team
	id := c.Param("id")

	if err := config.DB.First(&team, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tim tidak ditemukan!"})
		return
	}

	config.DB.Delete(&team)
	c.JSON(http.StatusOK, gin.H{"message": "Tim berhasil dihapus!"})
}

// AddPlayerToTeam mendaftarkan pemain ke dalam tim tertentu. Hanya admin.
// Endpoint: POST /api/teams/:id/players
func AddPlayerToTeam(c *gin.Context) {
	teamID := c.Param("id")

	// Cek dulu apakah tim-nya ada
	var team models.Team
	if err := config.DB.First(&team, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tim tidak ditemukan!"})
		return
	}

	// Bind input pemain baru
	var input models.Player
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak sesuai!"})
		return
	}

	// Paksa TeamID dari URL param, bukan dari body
	// Ini mencegah user mengisi team_id sembarangan
	input.TeamID = team.ID

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan pemain"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pemain berhasil ditambahkan ke tim!",
		"data":    input,
	})
}
