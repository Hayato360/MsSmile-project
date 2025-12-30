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

func TestCreateFetalKickCount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	r := gin.Default()
	r.POST("/kick-counts", controller.CreateFetalKickCount)

	// Setup: Create Pregnant Woman and Active Pregnancy
	db := config.DB()
	newMum := entity.PregnantWoman{Username: "KickMum", Password: "123"}
	db.Create(&newMum)

	preg := entity.Pregnancy{
		PregnantWomanID: &newMum.ID,
		Status:          "Active",
		LMP:             time.Now().AddDate(0, -5, 0),
		EDC:             time.Now().AddDate(0, 4, 0),
	}
	db.Create(&preg)

	// Case 1: Success Create
	payload := entity.FetalKickCount{
		PregnancyID:      &preg.ID,
		CountDate:        time.Now(),
		KickCountMorning: 10,
		KickCountLunch:   5,
		KickCountEvening: 8,
	}
	jsonValue, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/kick-counts", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)


	// Case 2: Update same day (Upsert)
	payloadUpdate := entity.FetalKickCount{
		PregnancyID:      &preg.ID,
		CountDate:        time.Now(), // Same day
		KickCountMorning: 12, // Changed
		KickCountLunch:   5,
		KickCountEvening: 8,
	}
	jsonValueUpdate, _ := json.Marshal(payloadUpdate)
	reqUpdate, _ := http.NewRequest("POST", "/kick-counts", bytes.NewBuffer(jsonValueUpdate))
	wUpdate := httptest.NewRecorder()
	r.ServeHTTP(wUpdate, reqUpdate)

	assert.Equal(t, http.StatusOK, wUpdate.Code)
	
	// Verify Update
	var rec entity.FetalKickCount
	db.Where("pregnancy_id = ?", preg.ID).First(&rec)
	assert.Equal(t, 12, rec.KickCountMorning)


	// Case 3: Inactive Pregnancy
	inactivePreg := entity.Pregnancy{
		PregnantWomanID: &newMum.ID,
		Status:          "Ended",
	}
	db.Create(&inactivePreg)
	payloadInactive := entity.FetalKickCount{
		PregnancyID: &inactivePreg.ID,
		CountDate:   time.Now(),
	}
	jsonValueInactive, _ := json.Marshal(payloadInactive)
	reqInactive, _ := http.NewRequest("POST", "/kick-counts", bytes.NewBuffer(jsonValueInactive))
	wInactive := httptest.NewRecorder()
	r.ServeHTTP(wInactive, reqInactive)

	assert.Equal(t, http.StatusBadRequest, wInactive.Code)
}

func TestGetFetalKickCounts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	r := gin.Default()
	r.GET("/kick-counts/pregnancy/:id", controller.GetFetalKickCountsByPregnancyID)

	db := config.DB()
	preg := entity.Pregnancy{Status: "Active"}
	db.Create(&preg)

	// Add records
	rec1 := entity.FetalKickCount{PregnancyID: &preg.ID, CountDate: time.Now().Add(-24 * time.Hour), KickCountMorning: 5}
	rec2 := entity.FetalKickCount{PregnancyID: &preg.ID, CountDate: time.Now(), KickCountMorning: 10}
	db.Create(&rec1)
	db.Create(&rec2)

	url := fmt.Sprintf("/kick-counts/pregnancy/%d", preg.ID)
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var data []entity.FetalKickCount
	json.Unmarshal(w.Body.Bytes(), &data)
	assert.Len(t, data, 2)
	assert.Equal(t, 5, data[0].KickCountMorning)
}
