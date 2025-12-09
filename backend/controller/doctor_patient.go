package controller

import (
	"net/http"
	"time"

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
		Preload("Appointment").
		First(&patient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	c.JSON(http.StatusOK, patient)
}

// POST /doctor/patient/:id/appointment
type CreateAppointmentInput struct {
	AppointmentDate string `json:"appointment_date"`
	Title           string `json:"title"`
	Location        string `json:"location"`
}

func DoctorCreateAppointment(c *gin.Context) {
	id := c.Param("id")
	var input CreateAppointmentInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()

	// Parse date
	// Expected format: "2025-11-25T09:00:00Z" or similar ISO
	date, err := time.Parse(time.RFC3339, input.AppointmentDate)
	if err != nil {
		// Try simple date format if ISO fails
		date, err = time.Parse("2006-01-02", input.AppointmentDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
	}

	// Create Appointment
	appt := entity.Appointment{
		AppointmentDate: date,
		Title:           input.Title,
		Location:        input.Location,
	}

	if err := db.Create(&appt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Link to Patient
	var patient entity.PregnantWoman
	if err := db.First(&patient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	patient.AppointmentID = &appt.ID
	if err := db.Save(&patient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update patient appointment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment created", "data": appt})
}
