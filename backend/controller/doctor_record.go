package controller

import (
	"net/http"

	"github.com/bestiesmile1845/Projecteiei/config"
	"github.com/bestiesmile1845/Projecteiei/entity"
	"github.com/gin-gonic/gin"
)

// POST /doctor/antenatal-visit - Doctor creates new visit record
func DoctorCreateAntenatalVisit(c *gin.Context) {
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

	c.JSON(http.StatusCreated, gin.H{"message": "Visit record created", "data": visit})
}

// GET /doctor/patient/:patientId/visits - Get all visits for a patient
func GetPatientVisits(c *gin.Context) {
	patientId := c.Param("patientId")
	var patient entity.PregnantWoman

	db := config.DB()

	// First get patient to find their pregnancy
	if err := db.Preload("Pregnancies").First(&patient, patientId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	// Get visits for patient's pregnancy (assuming latest pregnancy)
	if len(patient.Pregnancies) == 0 {
		c.JSON(http.StatusOK, []entity.AntenatalVisit{})
		return
	}

	var visits []entity.AntenatalVisit
	pregnancyId := patient.Pregnancies[len(patient.Pregnancies)-1].ID

	if err := db.Where("pregnancy_id = ?", pregnancyId).Find(&visits).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, visits)
}
