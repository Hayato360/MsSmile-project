package controller

import (
	"net/http"

	"github.com/bestiesmile1845/Projecteiei/config"
	"github.com/bestiesmile1845/Projecteiei/entity"
	"github.com/gin-gonic/gin"
)

// GET /vaccinations/pregnant-woman/:id
func GetVaccinationsByPregnantWomanID(c *gin.Context) {
	id := c.Param("id")
	var vaccinations []entity.Vaccination

	db := config.DB()

	if err := db.Preload("VaccineType").Preload("VacDose").Where("p_id = ?", id).Find(&vaccinations).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vaccinations)
}

// PUT /vaccinations/:id
func UpdateVaccination(c *gin.Context) {
	id := c.Param("id")
	var vaccination entity.Vaccination

	db := config.DB()

	if err := db.First(&vaccination, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vaccination not found"})
		return
	}

	if err := c.ShouldBindJSON(&vaccination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&vaccination).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated success", "data": vaccination})
}
