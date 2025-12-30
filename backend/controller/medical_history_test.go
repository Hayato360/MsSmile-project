package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"
	"time"

	"github.com/bestiesmile1845/Projecteiei/config"
	"github.com/bestiesmile1845/Projecteiei/controller"
	"github.com/bestiesmile1845/Projecteiei/entity"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDoctorCreateMedicalHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	r := gin.Default()
	r.POST("/doctor/medical-history", controller.DoctorCreateMedicalHistory)

	// Create a patient
	db := config.DB()
	newMum := entity.PregnantWoman{
		Username: "HistoryMum",
		Password: "123",
		BirthDate: time.Now(),
	}
	db.Create(&newMum)

	// Case 1: Create New
	payload := entity.MedicalHistory{
		PregnantWomanID: &newMum.ID,
		ChronicDiseases: "None",
		DrugAllergies:   "Penicillin",
	}
	jsonValue, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/doctor/medical-history", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Created", resp["message"])

	// Case 2: Update Existing (Upsert logic)
	payloadUpdate := entity.MedicalHistory{
		PregnantWomanID: &newMum.ID,
		ChronicDiseases: "Minor",
		DrugAllergies:   "Penicillin",
	}
	jsonValueUpdate, _ := json.Marshal(payloadUpdate)
	reqUpdate, _ := http.NewRequest("POST", "/doctor/medical-history", bytes.NewBuffer(jsonValueUpdate))
	wUpdate := httptest.NewRecorder()
	r.ServeHTTP(wUpdate, reqUpdate)

	assert.Equal(t, http.StatusOK, wUpdate.Code)
	json.Unmarshal(wUpdate.Body.Bytes(), &resp)
	assert.Equal(t, "Updated", resp["message"])
	
	// Verify in DB
	var mh entity.MedicalHistory
	db.Where("pregnant_woman_id = ?", newMum.ID).First(&mh)
	assert.Equal(t, "Minor", mh.ChronicDiseases)
}

func TestGetMedicalHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	r := gin.Default()
	r.GET("/medical-histories/pregnant-woman/:id", controller.GetMedicalHistoryByPregnantWomanID)

	// Create a patient & history
	db := config.DB()
	newMum := entity.PregnantWoman{
		Username: "GetHistoryMum",
		Password: "123",
		BirthDate: time.Now(),
	}
	db.Create(&newMum)
	
	mh := entity.MedicalHistory{
		PregnantWomanID: &newMum.ID,
		ChronicDiseases: "None",
	}
	db.Create(&mh)

	url := fmt.Sprintf("/medical-histories/pregnant-woman/%d", newMum.ID)
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "None", data["ChronicDiseases"]) // Field is Capitalized in struct, usually json is snake_case or whatever tag says.
	// But wait, struct definition didn't show tags for ChronicDiseases!
	// Struct: ChronicDiseases string
	// Default GORM/JSON: "ChronicDiseases" (if no tag).
	// Let's check struct again.
}
