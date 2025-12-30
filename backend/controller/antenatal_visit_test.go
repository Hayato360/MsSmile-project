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

func TestCreateAntenatalVisit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	r := gin.Default()
	r.POST("/antenatal-visits", controller.CreateAntenatalVisit)

	db := config.DB()
	newMum := entity.PregnantWoman{Username: "VisitMum", Password: "123"}
	db.Create(&newMum)

	preg := entity.Pregnancy{
		PregnantWomanID: &newMum.ID,
		Status:          "Active",
	}
	db.Create(&preg)

	// Case 1: Create Success
	payload := entity.AntenatalVisit{
		PregnancyID:   &preg.ID,
		VisitDate:     time.Now(),
		Weight:        60.5,
		BloodPressure: "120/80",
	}
	jsonValue, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/antenatal-visits", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetAntenatalVisits(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	r := gin.Default()
	r.GET("/antenatal-visits/pregnancy/:id", controller.GetAntenatalVisitsByPregnancyID)

	db := config.DB()
	preg := entity.Pregnancy{Status: "Active"}
	db.Create(&preg)

	visit := entity.AntenatalVisit{
		PregnancyID: &preg.ID,
		Weight:      55.0,
	}
	db.Create(&visit)

	url := fmt.Sprintf("/antenatal-visits/pregnancy/%d", preg.ID)
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var res []entity.AntenatalVisit
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.NotEmpty(t, res)
	assert.Equal(t, 55.0, res[0].Weight)
}
