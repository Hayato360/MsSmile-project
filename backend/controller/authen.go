package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bestiesmile1845/Projecteiei/config"
	"github.com/bestiesmile1845/Projecteiei/entity"  // ตรวจสอบเส้นทางให้ถูกต้อง
	"github.com/bestiesmile1845/Projecteiei/service" // ตรวจสอบเส้นทางให้ถูกต้อง
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// LoginPayload รับค่าจาก frontend
type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse สำหรับตอบกลับ frontend
type LoginResponse struct {
	Token string      `json:"token"`
	ID    uint        `json:"id"`
	User  interface{} `json:"user"` // สามารถเป็น entity.PregnantWoman หรือ entity.Doctor
	Role  string      `json:"role"` // Role ที่กำหนดโดยตรง
	Name  string      `json:"name"`
}

// --- Helper struct to hold common login data ---
type UserInterface struct {
	ID       uint
	Username string
	Password string
}

// POST /login
func Login(c *gin.Context) {
	var payload LoginPayload
	var user UserInterface
	var found bool = false
	var role string // **กำหนด Role โดยตรง**
	var userID uint
	var userName string
	var userDetails interface{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB() // ใช้ config.DB() แทน

	// 1. ลองค้นหาในตาราง Doctor
	var doc entity.Doctor
	if err := db.
		Where("username = ? OR email = ?", payload.Username, payload.Username).
		First(&doc).Error; err == nil {

		user = UserInterface{ID: doc.ID, Username: doc.Username, Password: doc.Password}
		userDetails = doc
		role = "doctor" // **กำหนด Role เป็น doctor**
		found = true
	}

	// 2. ถ้าไม่พบใน Doctor ลองค้นหาในตาราง PregnantWoman
	if !found {
		var woman entity.PregnantWoman
		if err := db.
			Where("username = ? OR email = ?", payload.Username, payload.Username).
			First(&woman).Error; err == nil {

			user = UserInterface{ID: woman.ID, Username: woman.Username, Password: woman.Password}
			userDetails = woman
			role = "pregnant" // **กำหนด Role เป็น pregnant**
			found = true
		}
	}

	// ถ้าไม่พบผู้ใช้ในทั้งสองตาราง
	if !found {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// ตรวจสอบรหัสผ่าน
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// สร้าง JWT token
	jwtWrapper := service.JwtWrapper{
		SecretKey:       "SvNQpBN8y3qlVrsGAYYWoJJk56LtzFHx",
		Issuer:          "AuthService",
		ExpirationHours: 24,
	}
	signedToken, err := jwtWrapper.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error signing token"})
		return
	}

	// Set final response details based on the user object (doc or woman)
	switch details := userDetails.(type) {
	case entity.PregnantWoman:
		userID = details.ID
		userName = details.FullName
	case entity.Doctor:
		userID = details.ID
		userName = details.FullName
	default:
		// ไม่ควรเกิดขึ้น
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unknown user type"})
		return
	}

	// DEBUG: แสดง role ที่ดึงมา
	fmt.Println("Role determined:", role)

	response := LoginResponse{
		Token: signedToken,
		ID:    userID,
		User:  userDetails,
		Role:  role, // ใช้ role ที่กำหนดโดยตรง
		Name:  userName,
	}

	// ส่ง response กลับในรูปแบบ {"data": response}
	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GET /me
func GetMe(c *gin.Context) {
	// ดึง Authorization header
	clientToken := c.Request.Header.Get("Authorization")
	if clientToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header provided"})
		return
	}

	// แยก Bearer token
	extractedToken := strings.Split(clientToken, "Bearer ")
	if len(extractedToken) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect Format of Authorization Token"})
		return
	}
	clientToken = strings.TrimSpace(extractedToken[1])

	// Validate token
	jwtWrapper := service.JwtWrapper{
		SecretKey: "SvNQpBN8y3qlVrsGAYYWoJJk56LtzFHx",
		Issuer:    "AuthService",
	}

	claims, err := jwtWrapper.ValidateToken(clientToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
		return
	}

	// ใช้ username จาก token claims เพื่อค้นหาผู้ใช้
	username := claims.Username
	db := config.DB()

	var found bool = false
	var role string
	var userID uint
	var userName string
	var userDetails interface{}

	// ค้นหาใน Doctor
	var doc entity.Doctor
	if err := db.Where("username = ?", username).First(&doc).Error; err == nil {
		userDetails = doc
		role = "doctor"
		userID = doc.ID
		userName = doc.FullName
		found = true
	}

	// ถ้าไม่พบใน Doctor ค้นหาใน PregnantWoman
	if !found {
		var woman entity.PregnantWoman
		if err := db.Preload("Pregnancies").Preload("Husband").Preload("MedicalHistories").Preload("Appointment").Where("username = ?", username).First(&woman).Error; err == nil {
			userDetails = woman
			role = "pregnant"
			userID = woman.ID
			userName = woman.FullName
			found = true
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// ส่งข้อมูลผู้ใช้กลับ
	response := gin.H{
		"id":   userID,
		"role": role,
		"name": userName,
		"user": userDetails,
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// RegisterDoctorPayload สำหรับลงทะเบียนแพทย์
type RegisterDoctorPayload struct {
	// ข้อมูลส่วนตัว
	FullName        string `json:"full_name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	PhoneNumber     string `json:"phone_number" binding:"required"`
	
	// ข้อมูลยืนยันตัวตน
	CitizenID       string `json:"citizen_id" binding:"required"`
	HospitalCode    string `json:"hospital_code" binding:"required"`
	DoctorLicenseNo string `json:"doctor_license_no" binding:"required"`

	// ข้อมูลบัญชี
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// POST /register/doctor
func RegisterDoctor(c *gin.Context) {
	var payload RegisterDoctorPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()

	// Step 1: ตรวจสอบว่า CitizenID ซ้ำหรือไม่
	var existingByCitizenID entity.Doctor
	if err := db.Where("citizen_id = ?", payload.CitizenID).First(&existingByCitizenID).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "เลขบัตรประชาชนนี้ถูกใช้งานแล้ว",
		})
		return
	}

	// Step 2: ตรวจสอบว่า username ซ้ำหรือไม่
	var existingByUsername entity.Doctor
	if err := db.Where("username = ?", payload.Username).First(&existingByUsername).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "ชื่อผู้ใช้นี้ถูกใช้งานแล้ว",
		})
		return
	}

	// Step 3: เข้ารหัสรหัสผ่าน
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "เกิดข้อผิดพลาดในการเข้ารหัสรหัสผ่าน",
		})
		return
	}

	// Step 4: สร้างบัญชีแพทย์ใหม่
	doctor := entity.Doctor{
		Username:        payload.Username,
		Password:        string(hashedPassword),
		CitizenID:       payload.CitizenID,
		HospitalCode:    payload.HospitalCode,
		DoctorLicenseNo: payload.DoctorLicenseNo,
		FullName:        payload.FullName,
		Email:           payload.Email,
		PhoneNumber:     payload.PhoneNumber,
		IsRegistered:    true,
	}

	if err := db.Create(&doctor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "เกิดข้อผิดพลาดในการลงทะเบียน",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "ลงทะเบียนสำเร็จ",
		"data": gin.H{
			"id":       doctor.ID,
			"username": doctor.Username,
			"name":     doctor.FullName,
		},
	})
}

