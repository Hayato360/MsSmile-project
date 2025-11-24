package controller

import (
	"net/http"

	"github.com/bestiesmile1845/Projecteiei/config"
	"github.com/bestiesmile1845/Projecteiei/entity"
	"github.com/gin-gonic/gin"
)

// GET /medical-histories/pregnant-woman/:id
func GetMedicalHistoryByPregnantWomanID(c *gin.Context) {
	id := c.Param("id")
	var history entity.MedicalHistory

	db := config.DB()

	if err := db.Where("pregnant_woman_id = ?", id).First(&history).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

// PUT /medical-histories/:id
func UpdateMedicalHistory(c *gin.Context) {
	id := c.Param("id")
	var history entity.MedicalHistory

	db := config.DB()

	if err := db.First(&history, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Medical History not found"})
		return
	}

	if err := c.ShouldBindJSON(&history); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated success", "data": history})
}
