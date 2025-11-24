package controller

import (
	"net/http"

	"github.com/bestiesmile1845/Projecteiei/config"
	"github.com/bestiesmile1845/Projecteiei/entity"
	"github.com/gin-gonic/gin"
)

// POST /kick-counts
func CreateFetalKickCount(c *gin.Context) {
	var kickCount entity.FetalKickCount

	// Bind JSON payload
	if err := c.ShouldBindJSON(&kickCount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()

	// Validate PregnancyID
	if kickCount.PregnancyID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PregnancyID is required"})
		return
	}

	// Create record
	if err := db.Create(&kickCount).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Created success", "data": kickCount})
}

// GET /kick-counts/pregnancy/:id
func GetFetalKickCountsByPregnancyID(c *gin.Context) {
	id := c.Param("id")
	var kickCounts []entity.FetalKickCount

	db := config.DB()

	// Find records by PregnancyID
	if err := db.Where("pregnancy_id = ?", id).Find(&kickCounts).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kickCounts)
}
