package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"laga-liga-backend/config"
	"laga-liga-backend/models"
)

// GetSports mengambil semua kategori olahraga.
// Endpoint: GET /api/sports
func GetSports(c *gin.Context) {
	var sports []models.Sport
	config.DB.Order("name ASC").Find(&sports)

	c.JSON(http.StatusOK, gin.H{
		"total": len(sports),
		"data":  sports,
	})
}
