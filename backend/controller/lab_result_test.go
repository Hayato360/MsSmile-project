package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
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

func TestDoctorCreateLabResult(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	r := gin.Default()
	r.POST("/doctor/lab-result", controller.DoctorCreateLabResult)

	// Create Directories for upload (mocking FS)
	os.MkdirAll("uploads/lab_results", os.ModePerm)
	defer os.RemoveAll("uploads") // Cleanup

	db := config.DB()
	preg := entity.Pregnancy{
		Status: "Active",
	}
	db.Create(&preg)

	// Prepare Multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add Fields
	writer.WriteField("PregnancyID", fmt.Sprintf("%d", preg.ID))
	writer.WriteField("TestDate", "2023-01-01")
	writer.WriteField("Hct", "35.5")
	writer.WriteField("Hb", "12.2")
	writer.WriteField("HbTyping", "Normal")
	writer.WriteField("OtherRemarks", "None")

	// Add File
	part, _ := writer.CreateFormFile("File", "test_report.pdf")
	part.Write([]byte("dummy pdf content"))

	writer.Close()

	req, _ := http.NewRequest("POST", "/doctor/lab-result", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Created", resp["message"])
	
	// Verify DB
	var lr entity.LabResult
	db.First(&lr)
	assert.Equal(t, 35.5, lr.Hct)
	assert.Contains(t, lr.FilePath, "test_report.pdf")
}

func TestGetLabResults(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	r := gin.Default()
	r.GET("/doctor/pregnancy/:pregnancyId/lab-results", controller.GetLabResultsByPregnancyID)

	db := config.DB()
	preg := entity.Pregnancy{Status: "Active"}
	db.Create(&preg)

	lr := entity.LabResult{
		PregnancyID: &preg.ID,
		Hct:         30.0,
		TestDate:    time.Now(),
	}
	db.Create(&lr)

	url := fmt.Sprintf("/doctor/pregnancy/%d/lab-results", preg.ID)
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var data map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &data)
	// data is {"data": [...]}
	results := data["data"].([]interface{})
	assert.NotEmpty(t, results)
	assert.Equal(t, 30.0, results[0].(map[string]interface{})["Hct"])
}
