package controller

import (
	"net/http"
	"time"

	"github.com/bestiesmile1845/Projecteiei/config"
	"github.com/bestiesmile1845/Projecteiei/entity"
	"github.com/gin-gonic/gin"
)

// PUT /profile/husband - Update husband information
func UpdateHusband(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID := val.(uint)
	var input struct {
		FullName    string `json:"full_name"`
		Age         int    `json:"age"`
		CitizenID   string `json:"citizen_id"`
		PhoneNumber string `json:"phone_number"`
		Email       string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()

	// Get user
	var user entity.PregnantWoman
	if err := db.Preload("Husband").First(&user, uint(userID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Create or update husband
	if user.HusbandID == nil {
		tx := db.Begin()

		// Create new husband
		husband := entity.Husband{
			FullName:    input.FullName,
			Age:         input.Age,
			CitizenID:   input.CitizenID,
			PhoneNumber: input.PhoneNumber,
			Email:       input.Email,
		}
		if err := tx.Create(&husband).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Update user with husband ID
		if err := tx.Model(&user).Update("husband_id", husband.ID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"message": "Husband created", "data": husband})
	} else {
		// Update existing husband
		husband := user.Husband
		husband.FullName = input.FullName
		husband.Age = input.Age
		husband.CitizenID = input.CitizenID
		husband.PhoneNumber = input.PhoneNumber
		husband.Email = input.Email

		if err := db.Save(husband).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Husband updated", "data": husband})
	}
}

// PUT /profile/personal - Update personal information
func UpdatePersonalProfile(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID := val.(uint)

	var input struct {
		FullName    string `json:"full_name"`
		BirthDate   string `json:"birth_date"` // Receive as string to parse
		CitizenID   string `json:"citizen_id"`
		PhoneNumber string `json:"phone_number"`
		Email       string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()
	var user entity.PregnantWoman
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.FullName = input.FullName
	user.CitizenID = input.CitizenID
	user.PhoneNumber = input.PhoneNumber
	user.Email = input.Email

	// Parse BirthDate
	if input.BirthDate != "" {
		parsedDate, err := time.Parse("2006-01-02", input.BirthDate) // Assuming YYYY-MM-DD
		if err == nil {
			user.BirthDate = parsedDate
			// Recalculate Age
			now := time.Now()
			age := now.Year() - parsedDate.Year()
			if now.YearDay() < parsedDate.YearDay() {
				age--
			}
			user.Age = age
		} else {
			// Try RFC3339 if simple date fails
			parsedDate, err := time.Parse(time.RFC3339, input.BirthDate)
			if err == nil {
				user.BirthDate = parsedDate
				// Recalculate Age
				now := time.Now()
				age := now.Year() - parsedDate.Year()
				if now.YearDay() < parsedDate.YearDay() {
					age--
				}
				user.Age = age
			}
		}
	}

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Personal info updated", "data": user})
}

// PUT /profile/medical-history - Update or Create Medical History
func UpdateMyMedicalHistory(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID := val.(uint)

	var input entity.MedicalHistory
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()
	var user entity.PregnantWoman
	// Preload MedicalHistories to check if exists
	if err := db.Preload("MedicalHistories").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if len(user.MedicalHistories) > 0 {
		// Update existing
		history := user.MedicalHistories[0] // Assume first one
		// Update fields
		history.ChronicDiseases = input.ChronicDiseases
		history.HeartDisease = input.HeartDisease
		history.Thyroid = input.Thyroid
		history.OtherDiseases = input.OtherDiseases
		history.SurgeryHistory = input.SurgeryHistory
		history.OtherSurgery = input.OtherSurgery
		history.GeneticDiseases = input.GeneticDiseases
		history.DrugAllergies = input.DrugAllergies
		history.FamilyHistoryHT = input.FamilyHistoryHT
		history.FamilyHistoryDiabetes = input.FamilyHistoryDiabetes
		history.FamilyHistoryThalassemia = input.FamilyHistoryThalassemia
		history.FamilyHistoryCongenital = input.FamilyHistoryCongenital
		history.OtherFamilyHistory = input.OtherFamilyHistory
		history.ContraceptionBeforeMethod = input.ContraceptionBeforeMethod
		history.ContraceptionBeforeDuration = input.ContraceptionBeforeDuration
		history.ContraceptionLastMethod = input.ContraceptionLastMethod
		history.ContraceptionLastDuration = input.ContraceptionLastDuration
		history.MenstrualCycle = input.MenstrualCycle
		history.MenstrualDuration = input.MenstrualDuration
		history.MenstrualCondition = input.MenstrualCondition

		if err := db.Save(&history).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Medical history updated", "data": history})
	} else {
		// Create new
		input.ID = 0 // Reset ID to ensure creation
		input.PregnantWomanID = &userID
		if err := db.Create(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Medical history created", "data": input})
	}
}
