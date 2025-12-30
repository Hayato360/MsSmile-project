package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bestiesmile1845/Projecteiei/config"
	"github.com/bestiesmile1845/Projecteiei/controller"
	"github.com/bestiesmile1845/Projecteiei/entity"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDoctorCreateVaccination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	r := gin.Default()
	r.POST("/doctor/vaccination", controller.DoctorCreateVaccination)

	// Create Patient
	db := config.DB()
	newMum := entity.PregnantWoman{Username: "VacMum", Password: "123"}
	db.Create(&newMum)

	// Create Vaccine Type
	vacType := entity.VaccineType{Name: "COVID-19"}
	db.Create(&vacType)

	// Helper for time pointer
	date2023 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	
	// Case 1: Create New
	payload := entity.Vaccination{
		PregnantWomanID:       &newMum.ID, 
		VaccineTypeID:         &vacType.ID,
		PreviousDoses:         2,
		LastPreviousDateYear:  &date2023,
	}
	// Note: Check entity definition for exact field names/pointers if previous tests failed on structure.
	// But based on controller `payload.PregnantWomanID`, it seems it might match `MedicalHistory` pattern.
	
	jsonValue, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/doctor/vaccination", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	
	// Case 2: Update Existing (Upsert)
	date2024 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	payloadUpdate := entity.Vaccination{
		PregnantWomanID:       &newMum.ID,
		VaccineTypeID:         &vacType.ID,
		PreviousDoses:         3, // Changed
		LastPreviousDateYear:  &date2024,
	}
	jsonValueUpdate, _ := json.Marshal(payloadUpdate)
	reqUpdate, _ := http.NewRequest("POST", "/doctor/vaccination", bytes.NewBuffer(jsonValueUpdate))
	wUpdate := httptest.NewRecorder()
	r.ServeHTTP(wUpdate, reqUpdate)

	assert.Equal(t, http.StatusOK, wUpdate.Code)
	
	// Verify Update
	var vac entity.Vaccination
	db.Where("p_id = ?", newMum.ID).First(&vac)
	assert.Equal(t, 3, vac.PreviousDoses)
}

func TestUpdateVaccination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	r := gin.Default()
	r.PUT("/vaccinations/:id", controller.UpdateVaccination)

	db := config.DB()
	vac := entity.Vaccination{PreviousDoses: 1}
	db.Create(&vac)

	payload := map[string]interface{}{
		"PreviousDoses": 5,
	}
	jsonValue, _ := json.Marshal(payload)
	url := fmt.Sprintf("/vaccinations/%d", vac.ID)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var updatedVac entity.Vaccination
	db.First(&updatedVac, vac.ID)
	assert.Equal(t, 5, updatedVac.PreviousDoses)
}

func TestGetVaccinations(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	r := gin.Default()
	r.GET("/vaccinations/pregnant-woman/:id", controller.GetVaccinationsByPregnantWomanID)

	db := config.DB()
	newMum := entity.PregnantWoman{Username: "GetVacMum", Password: "123"}
	db.Create(&newMum)
	
	vacType := entity.VaccineType{Name: "Flu"}
	db.Create(&vacType)

	vac := entity.Vaccination{
		PregnantWomanID: &newMum.ID,
		VaccineTypeID:   &vacType.ID,
		PreviousDoses:   1,
	}
	db.Create(&vac)

	url := fmt.Sprintf("/vaccinations/pregnant-woman/%d", newMum.ID)
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var res []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.NotEmpty(t, res)
	// Verify Preload (VaccineType) - if controller preloads it
	// assert.Equal(t, "Flu", res[0]["VaccineType"].(map[string]interface{})["Name"])
	// Check field casing logic
}

func TestListVaccineTypes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	r := gin.Default()
	r.GET("/vaccine-types", controller.ListVaccineTypes)

	db := config.DB()
	db.Create(&entity.VaccineType{Name: "TypeA"})
	db.Create(&entity.VaccineType{Name: "TypeB"})

	req, _ := http.NewRequest("GET", "/vaccine-types", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	data := res["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 2)
}
