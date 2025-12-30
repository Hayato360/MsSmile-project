package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bestiesmile1845/Projecteiei/config"
	"github.com/bestiesmile1845/Projecteiei/controller"
	"github.com/bestiesmile1845/Projecteiei/entity"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock Auth Middleware to inject UserID
func mockAuthMiddleware(userID uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	}
}

func TestUpdateHusband(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	// Setup User
	db := config.DB()
	user := entity.PregnantWoman{
		Username: "HusbandTestUser",
		Password: "123",
	}
	db.Create(&user)

	r := gin.Default()
	// Use mock middleware instead of real JWT parsing
	protected := r.Group("/")
	protected.Use(mockAuthMiddleware(user.ID))
	protected.PUT("/profile/husband", controller.UpdateHusband)

	// Case 1: Create Husband
	payload := map[string]interface{}{
		"full_name":    "Husband A",
		"age":          30,
		"citizen_id":   "1234567890123",
		"phone_number": "0812345678",
		"email":        "husband@mail.com",
	}
	jsonValue, _ := json.Marshal(payload)
	
	req, _ := http.NewRequest("PUT", "/profile/husband", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	// Verify Linkage
	var updatedUser entity.PregnantWoman
	db.Preload("Husband").First(&updatedUser, user.ID)
	assert.NotNil(t, updatedUser.HusbandID)
	assert.Equal(t, "Husband A", updatedUser.Husband.FullName)

	// Case 2: Update Husband
	payload["full_name"] = "Husband B"
	jsonValueUpdate, _ := json.Marshal(payload)
	reqUpdate, _ := http.NewRequest("PUT", "/profile/husband", bytes.NewBuffer(jsonValueUpdate))
	wUpdate := httptest.NewRecorder()
	r.ServeHTTP(wUpdate, reqUpdate)

	assert.Equal(t, http.StatusOK, wUpdate.Code)

	db.Preload("Husband").First(&updatedUser, user.ID)
	assert.Equal(t, "Husband B", updatedUser.Husband.FullName)
}

func TestUpdatePersonalProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	// Setup User
	db := config.DB()
	user := entity.PregnantWoman{
		Username: "PersonalTestUser",
		Password: "123",
		FullName: "Old Name",
	}
	db.Create(&user)

	r := gin.Default()
	protected := r.Group("/")
	protected.Use(mockAuthMiddleware(user.ID))
	protected.PUT("/profile/personal", controller.UpdatePersonalProfile)

	// Update
	payload := map[string]interface{}{
		"full_name":    "New Name",
		"birth_date":   "2000-01-01",
		"citizen_id":   "9999999999999",
		"phone_number": "0999999999",
		"email":        "new@mail.com",
	}
	jsonValue, _ := json.Marshal(payload)
	
	req, _ := http.NewRequest("PUT", "/profile/personal", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var updatedUser entity.PregnantWoman
	db.First(&updatedUser, user.ID)
	assert.Equal(t, "New Name", updatedUser.FullName)
	assert.Equal(t, "0999999999", updatedUser.PhoneNumber)
	// Check Age calculation (2000 -> 25 in 2025)
	// Note: Verify logic in controller for age
}

func TestUpdateMyMedicalHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	// Setup User
	db := config.DB()
	user := entity.PregnantWoman{
		Username: "MyHistoryTestUser",
		Password: "123",
	}
	db.Create(&user)

	r := gin.Default()
	protected := r.Group("/")
	protected.Use(mockAuthMiddleware(user.ID))
	protected.PUT("/profile/medical-history", controller.UpdateMyMedicalHistory)

	// Create
	payload := entity.MedicalHistory{
		ChronicDiseases: "Asthma",
	}
	jsonValue, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", "/profile/medical-history", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	// Verify
	var histories []entity.MedicalHistory
	db.Where("pregnant_woman_id = ?", user.ID).Find(&histories)
	assert.NotEmpty(t, histories)
	assert.Equal(t, "Asthma", histories[0].ChronicDiseases)
}
