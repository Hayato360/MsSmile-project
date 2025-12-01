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

// GET /vaccine-types
func ListVaccineTypes(c *gin.Context) {
	var vaccineTypes []entity.VaccineType
	db := config.DB()
	if err := db.Find(&vaccineTypes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": vaccineTypes})
}

// POST /doctor/vaccination
func DoctorCreateVaccination(c *gin.Context) {
	var payload entity.Vaccination
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()
	var existingVaccination entity.Vaccination

	// Check if exists
	if err := db.Where("p_id = ? AND vaccine_type_id = ?", payload.PregnantWomanID, payload.VaccineTypeID).First(&existingVaccination).Error; err == nil {
		// Exists -> Update
		existingVaccination.IsPreviouslyVaccinated = payload.IsPreviouslyVaccinated
		existingVaccination.PreviousDoses = payload.PreviousDoses
		existingVaccination.LastPreviousDateYear = payload.LastPreviousDateYear
		existingVaccination.Dose1DateDuringPreg = payload.Dose1DateDuringPreg
		existingVaccination.Dose2DateDuringPreg = payload.Dose2DateDuringPreg
		existingVaccination.Dose3DateDuringPreg = payload.Dose3DateDuringPreg
		existingVaccination.IsHistoryUnknown = payload.IsHistoryUnknown
		existingVaccination.ReasonForNotVaccinating = payload.ReasonForNotVaccinating
		existingVaccination.Remarks = payload.Remarks

		if err := db.Save(&existingVaccination).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Updated success", "data": existingVaccination})
		return
	}

	// Not exists -> Create
	if err := db.Create(&payload).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Created success", "data": payload})
}
