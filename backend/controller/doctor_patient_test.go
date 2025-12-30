package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/bestiesmile1845/Projecteiei/config"
	"github.com/bestiesmile1845/Projecteiei/controller"
	"github.com/bestiesmile1845/Projecteiei/entity"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)



func TestGetDoctorPatients(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	setupTestDB()

	r := gin.Default()
	r.GET("/doctor/patients", controller.GetDoctorPatients)

	// Execute
	req, _ := http.NewRequest("GET", "/doctor/patients", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var patients []entity.PregnantWoman
	err := json.Unmarshal(w.Body.Bytes(), &patients)
	assert.Nil(t, err)
	
	// Expect at least one patient from the seed data (Mommy)
	assert.NotEmpty(t, patients)
	assert.Equal(t, "Mommy", patients[0].Username)
}

func TestDoctorCreateAppointment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	r := gin.Default()
	r.POST("/doctor/patient/:id/appointment", controller.DoctorCreateAppointment)

	// Create a specific patient for this test to ensure we have a valid ID
	db := config.DB()
	newMum := entity.PregnantWoman{
		Username: "TestAppointmentMum",
		Password: "123",
		BirthDate: time.Now(),
	}
	db.Create(&newMum)
	id := newMum.ID

	payload := controller.CreateAppointmentInput{
		AppointmentDate: "2025-12-30T10:00:00Z",
		Title:           "Root Canal",
		Location:        "Clinic A",
	}
	jsonValue, _ := json.Marshal(payload)
	
	url := fmt.Sprintf("/doctor/patient/%d/appointment", id)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.Equal(t, "Appointment created", resp["message"])

	// Verify linkage
	var updatedMum entity.PregnantWoman
	db.First(&updatedMum, id)
	assert.NotNil(t, updatedMum.AppointmentID)
}
