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

func TestDoctorCreatePreviousPregnancy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	r := gin.Default()
	r.POST("/doctor/previous-pregnancy", controller.DoctorCreatePreviousPregnancy)

	db := config.DB()
	newMum := entity.PregnantWoman{Username: "PrevMum", Password: "123"}
	db.Create(&newMum)

	payload := entity.PreviousPregnancy{
		PregnantWomanID:       &newMum.ID,
		PregnancyNo:           1,
		DeliveryDate:          time.Now().AddDate(-2, 0, 0),
		Complications:         "None",
	}
	jsonValue, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/doctor/previous-pregnancy", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var rec entity.PreviousPregnancy
	db.First(&rec)
	assert.Equal(t, 1, rec.PregnancyNo)
}

func TestGetPreviousPregnancies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	r := gin.Default()
	r.GET("/doctor/patient/:patientId/previous-pregnancies", controller.GetPreviousPregnancies)

	db := config.DB()
	newMum := entity.PregnantWoman{Username: "GetPrevMum", Password: "123"}
	db.Create(&newMum)

	pp := entity.PreviousPregnancy{
		PregnantWomanID: &newMum.ID,
		PregnancyNo:     1,
		BirthWeight:     3000,
	}
	db.Create(&pp)

	url := fmt.Sprintf("/doctor/patient/%d/previous-pregnancies", newMum.ID)
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].([]interface{})
	assert.NotEmpty(t, data)
	assert.Equal(t, 3000.0, data[0].(map[string]interface{})["birth_weight"])
}

func TestDoctorCreatePregnancy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	r := gin.Default()
	r.POST("/doctor/pregnancy", controller.DoctorCreatePregnancy)

	db := config.DB()
	newMum := entity.PregnantWoman{Username: "PregMum", Password: "123"}
	db.Create(&newMum)

	payload := entity.Pregnancy{
		PregnantWomanID: &newMum.ID,
		PregnancyNo:     2,
	}
	jsonValue, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/doctor/pregnancy", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	
	var preg entity.Pregnancy
	db.Where("p_id = ?", newMum.ID).First(&preg)
	assert.Equal(t, 2, preg.PregnancyNo)
	assert.Equal(t, "Active", preg.Status) // Assuming default is active
}

func TestEndPregnancy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	r := gin.Default()
	r.POST("/doctor/pregnancy/:id/end", controller.EndPregnancy)

	db := config.DB()
	preg := entity.Pregnancy{Status: "Active"}
	db.Create(&preg)

	// Payload
	payload := map[string]interface{}{
		"delivery_date": time.Now(),
		"child_status": "Healthy",
	}
	jsonValue, _ := json.Marshal(payload)

	url := fmt.Sprintf("/doctor/pregnancy/%d/end", preg.ID)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var updatedPreg entity.Pregnancy
	db.First(&updatedPreg, preg.ID)
	assert.Equal(t, "Ended", updatedPreg.Status)
}
