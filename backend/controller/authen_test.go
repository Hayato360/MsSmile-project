package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bestiesmile1845/Projecteiei/controller"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	r := gin.Default()
	r.POST("/login", controller.Login)

	// Case 1: Success (Doctor)
	// Seeded in SetupDatabase: Username="Doctor", Password="123" (hashed)
	payload := controller.LoginPayload{
		Username: "Doctor",
		Password: "123456", // SetupDatabase uses "123456" for hash
	}
	jsonValue, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	// Check response structure
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	data, ok := resp["data"].(map[string]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, data["token"])
	assert.Equal(t, "doctor", data["role"])

	// Case 2: Success (PregnantWoman "Mommy")
	payloadMommy := controller.LoginPayload{
		Username: "Mommy",
		Password: "123456",
	}
	jsonValueMommy, _ := json.Marshal(payloadMommy)
	reqMommy, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValueMommy))
	wMommy := httptest.NewRecorder()
	r.ServeHTTP(wMommy, reqMommy)
	assert.Equal(t, http.StatusOK, wMommy.Code)


	// Case 3: Invalid Password
	payloadInvalid := controller.LoginPayload{
		Username: "Doctor",
		Password: "WrongPassword",
	}
	jsonValueInvalid, _ := json.Marshal(payloadInvalid)
	reqInvalid, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValueInvalid))
	wInvalid := httptest.NewRecorder()
	r.ServeHTTP(wInvalid, reqInvalid)

	assert.Equal(t, http.StatusUnauthorized, wInvalid.Code)
}

func TestGetMe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	r := gin.Default()
	r.GET("/me", controller.GetMe)

	// 1. Get a token first via Login
	loginPayload := controller.LoginPayload{
		Username: "Doctor",
		Password: "123456",
	}
	jsonValue, _ := json.Marshal(loginPayload)
	reqLogin, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue)) // Need to route logic or helpers
	// To avoid routing overhead, just reuse the logic or mock token.
	// But getting a real token is integration-style and safe.
	
	// Wait, controller.Login is not bound to a router here if I don't use the same 'r'.
	// Let's bind it just for helper usage.
	r.POST("/login", controller.Login)
	wLogin := httptest.NewRecorder()
	r.ServeHTTP(wLogin, reqLogin)
	var resp map[string]interface{}
	json.Unmarshal(wLogin.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	token := data["token"].(string)

	// 2. Use the token to GetMe
	req, _ := http.NewRequest("GET", "/me", nil)
	req.Header.Set("Authorization", "Bearer " + token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var meResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &meResp)
	meData := meResp["data"].(map[string]interface{})
	
	assert.Equal(t, "doctor", meData["role"])
	assert.Equal(t, "Doctor D", meData["name"])
}
