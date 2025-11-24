package controller

import (
	"net/http"

	"github.com/bestiesmile1845/Projecteiei/config"
	"github.com/bestiesmile1845/Projecteiei/entity"
	"github.com/gin-gonic/gin"
)

// POST /doctor/pregnancy - Create pregnancy record for patient
func DoctorCreatePregnancy(c *gin.Context) {
	var pregnancy entity.Pregnancy

	if err := c.ShouldBindJSON(&pregnancy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()

	if err := db.Create(&pregnancy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pregnancy record created", "data": pregnancy})
}
