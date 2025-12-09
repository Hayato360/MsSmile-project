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
		existing.HeartDisease = input.HeartDisease
		existing.Thyroid = input.Thyroid
		existing.OtherDiseases = input.OtherDiseases
		existing.SurgeryHistory = input.SurgeryHistory
		existing.OtherSurgery = input.OtherSurgery
		existing.GeneticDiseases = input.GeneticDiseases
		existing.DrugAllergies = input.DrugAllergies
		
		existing.FamilyHistoryHT = input.FamilyHistoryHT
		existing.FamilyHistoryDiabetes = input.FamilyHistoryDiabetes
		existing.FamilyHistoryThalassemia = input.FamilyHistoryThalassemia
		existing.FamilyHistoryCongenital = input.FamilyHistoryCongenital
		existing.OtherFamilyHistory = input.OtherFamilyHistory

		existing.ContraceptionBeforeMethod = input.ContraceptionBeforeMethod
		existing.ContraceptionBeforeDuration = input.ContraceptionBeforeDuration
		existing.ContraceptionLastMethod = input.ContraceptionLastMethod
		existing.ContraceptionLastDuration = input.ContraceptionLastDuration

		existing.MenstrualCycle = input.MenstrualCycle
		existing.MenstrualDuration = input.MenstrualDuration
		existing.MenstrualCondition = input.MenstrualCondition

		existing.Preeclampsia = input.Preeclampsia
		existing.Remarks = input.Remarks

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

// POST /doctor/previous-pregnancy - Create previous pregnancy record
func DoctorCreatePreviousPregnancy(c *gin.Context) {
	var input entity.PreviousPregnancy

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()

	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Created", "data": input})
}

// GET /doctor/patient/:patientId/previous-pregnancies
func GetPreviousPregnancies(c *gin.Context) {
	patientId := c.Param("patientId")
	db := config.DB()

	var pregnancies []entity.PreviousPregnancy
	if err := db.Where("pregnant_woman_id = ?", patientId).Order("pregnancy_no ASC").Find(&pregnancies).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"data": []entity.PreviousPregnancy{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": pregnancies})
}

// GET /previous-pregnancies/pregnant-woman/:id
func GetPreviousPregnanciesByPregnantWomanID(c *gin.Context) {
	id := c.Param("id")
	db := config.DB()

	var pregnancies []entity.PreviousPregnancy
	if err := db.Where("pregnant_woman_id = ?", id).Order("pregnancy_no ASC").Find(&pregnancies).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"data": []entity.PreviousPregnancy{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": pregnancies})
}

// POST /doctor/lab-result - Create lab result with file upload
func DoctorCreateLabResult(c *gin.Context) {
	// Parse multipart form (32MB limit)
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large or invalid multipart form"})
		return
	}

	// 1. Handle File Upload
	var filePath string
	file, err := c.FormFile("File")
	if err == nil {
		// Create uploads directory if it doesn't exist (handled by mkdir in most cases, but good to be safe)
		// Assuming "uploads/lab_results" structure
		filename := time.Now().Format("20060102150405") + "_" + file.Filename
		dst := "uploads/lab_results/" + filename
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file: " + err.Error()})
			return
		}
		filePath = dst
	}

	// 2. Bind other fields
	// Let's use a DTO struct for binding form data
	
	// Let's us a DTO struct inside for binding form data
	var formDto struct {
		PregnancyID  uint      `form:"PregnancyID"`
		TestDate     string    `form:"TestDate"` // Date coming as string "2023-01-01"
		Hct          float64   `form:"Hct"`
		Hb           float64   `form:"Hb"`
		HbTyping     string    `form:"HbTyping"`
		OtherRemarks string    `form:"OtherRemarks"`
	}

	if err := c.ShouldBind(&formDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse Date
	testDate, _ := time.Parse("2006-01-02", formDto.TestDate)

	db := config.DB()

	labResult := entity.LabResult{
		PregnancyID:  &formDto.PregnancyID,
		TestDate:     testDate,
		Hct:          formDto.Hct,
		Hb:           formDto.Hb,
		HbTyping:     formDto.HbTyping,
		OtherRemarks: formDto.OtherRemarks,
		FilePath:     filePath,
	}

	if err := db.Create(&labResult).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Created", "data": labResult})
}

// GET /doctor/pregnancy/:pregnancyId/lab-results
func GetLabResultsByPregnancyID(c *gin.Context) {
	id := c.Param("pregnancyId")
	db := config.DB()

	var results []entity.LabResult
	if err := db.Where("pregnancy_id = ?", id).Order("test_date DESC").Find(&results).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"data": []entity.LabResult{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
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
