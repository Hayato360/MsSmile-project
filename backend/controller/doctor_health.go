 package controller

import (
	"net/http"
	"time"

	"github.com/bestiesmile1845/Projecteiei/config"
	"github.com/bestiesmile1845/Projecteiei/entity"
	"github.com/gin-gonic/gin"
)

// POST /doctor/medical-history - Create or update medical history
func DoctorCreateMedicalHistory(c *gin.Context) {
	var input entity.MedicalHistory

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()

	// Check if medical history exists for the patient
	var existing entity.MedicalHistory
	result := db.Where("pregnant_woman_id = ?", input.PregnantWomanID).First(&existing)

	if result.Error == nil {
		// Update existing
		existing.ChronicDiseases = input.ChronicDiseases
		existing.GeneticDiseases = input.GeneticDiseases
		existing.DrugAllergies = input.DrugAllergies
		existing.FamilyHistoryHT = input.FamilyHistoryHT
		existing.FamilyHistoryCongenital = input.FamilyHistoryCongenital
		existing.OtherFamilyHistory = input.OtherFamilyHistory

		if err := db.Save(&existing).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Updated", "data": existing})
	} else {
		// Create new
		if err := db.Create(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Created", "data": input})
	}
}

// POST /doctor/lab-result - Create lab result
func DoctorCreateLabResult(c *gin.Context) {
	var input struct {
		PregnancyID  uint      `json:"PregnancyID"`
		TestDate     time.Time `json:"TestDate"`
		Hct          float64   `json:"Hct"`
		Hb           float64   `json:"Hb"`
		HbTyping     string    `json:"HbTyping"`
		OtherRemarks string    `json:"OtherRemarks"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()

	labResult := entity.LabResult{
		PregnancyID:  &input.PregnancyID,
		TestDate:     input.TestDate,
		Hct:          input.Hct,
		Hb:           input.Hb,
		HbTyping:     input.HbTyping,
		OtherRemarks: input.OtherRemarks,
	}

	if err := db.Create(&labResult).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Created", "data": labResult})
}

// POST /doctor/vaccination - Create vaccination record
func DoctorCreateVaccination(c *gin.Context) {
	var input struct {
		PregnantWomanID        uint       `json:"PregnantWomanID"`
		VaccineTypeID          uint       `json:"VaccineTypeID"`
		IsPreviouslyVaccinated bool       `json:"IsPreviouslyVaccinated"`
		PreviousDoses          int        `json:"PreviousDoses"`
		LastPreviousDateYear   *time.Time `json:"LastPreviousDateYear"`
		Dose1DateDuringPreg    *time.Time `json:"Dose1DateDuringPreg"`
		Dose2DateDuringPreg    *time.Time `json:"Dose2DateDuringPreg"`
		Remarks                string     `json:"Remarks"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()

	vaccination := entity.Vaccination{
		PregnantWomanID:        &input.PregnantWomanID,
		VaccineTypeID:          &input.VaccineTypeID,
		IsPreviouslyVaccinated: input.IsPreviouslyVaccinated,
		PreviousDoses:          input.PreviousDoses,
		LastPreviousDateYear:   input.LastPreviousDateYear,
		Dose1DateDuringPreg:    input.Dose1DateDuringPreg,
		Dose2DateDuringPreg:    input.Dose2DateDuringPreg,
		Remarks:                input.Remarks,
	}

	if err := db.Create(&vaccination).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Created", "data": vaccination})
}

// GET /doctor/patient/:patientId/medical-history
func GetPatientMedicalHistory(c *gin.Context) {
	patientId := c.Param("patientId")
	db := config.DB()

	var history entity.MedicalHistory
	if err := db.Where("pregnant_woman_id = ?", patientId).First(&history).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": history})
}
