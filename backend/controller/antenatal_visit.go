package controller

import (
	"net/http"

	"github.com/bestiesmile1845/Projecteiei/config"
	"github.com/bestiesmile1845/Projecteiei/entity"
	"github.com/gin-gonic/gin"
)

// GET /antenatal-visits/pregnancy/:id
func GetAntenatalVisitsByPregnancyID(c *gin.Context) {
	id := c.Param("id")
	var visits []entity.AntenatalVisit

	db := config.DB()

	if err := db.Where("pregnancy_id = ?", id).Find(&visits).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, visits)
}

// POST /antenatal-visits
func CreateAntenatalVisit(c *gin.Context) {
	var visit entity.AntenatalVisit

	if err := c.ShouldBindJSON(&visit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()

	if err := db.Create(&visit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Created success", "data": visit})
}
