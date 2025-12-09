package controller

import (
	"net/http"

	"github.com/bestiesmile1845/Projecteiei/config"
	"github.com/bestiesmile1845/Projecteiei/entity"
	"github.com/gin-gonic/gin"
)

// GET /doctor/patients - Get all patients (pregnant women)
func GetDoctorPatients(c *gin.Context) {
	var patients []entity.PregnantWoman

	db := config.DB()

	// For now, return all patients. In production, you'd filter by doctor assignment
	if err := db.Preload("Pregnancies").Preload("MedicalHistories").Find(&patients).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, patients)
}

// GET /doctor/patients/:id - Get specific patient details
func GetDoctorPatientDetail(c *gin.Context) {
	id := c.Param("id")
	var patient entity.PregnantWoman

	db := config.DB()

	if err := db.Preload("Pregnancies").
		Preload("MedicalHistories").
		Preload("Vaccinations").
		Preload("PregnancyHistories").
		First(&patient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	c.JSON(http.StatusOK, patient)
}
