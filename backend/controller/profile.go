package controller

import (
	"net/http"

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
