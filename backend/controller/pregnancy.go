package controller

import (
	"net/http"
	"time"

	"github.com/bestiesmile1845/Projecteiei/config"
	"github.com/bestiesmile1845/Projecteiei/entity"
	"github.com/gin-gonic/gin"
)

// POST /doctor/pregnancy - Create pregnancy record for patient
func DoctorCreatePregnancy(c *gin.Context) {
	var pregnancy entity.Pregnancy

	if err := c.ShouldBindJSON(&pregnancy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()

	// Check if there is already an active pregnancy for this patient
	var activePregnancy entity.Pregnancy
	if err := db.Where("pregnant_woman_id = ? AND status = ?", pregnancy.PregnantWomanID, "Active").First(&activePregnancy).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Patient already has an active pregnancy. Please end the current pregnancy first."})
		return
	}

	// Set default status to Active
	pregnancy.Status = "Active"

	// Calculate EDC if not provided (LMP + 280 days)
	if pregnancy.EDC.IsZero() && !pregnancy.LMP.IsZero() {
		pregnancy.EDC = pregnancy.LMP.AddDate(0, 0, 280)
	}

	if err := db.Create(&pregnancy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pregnancy record created", "data": pregnancy})
}

// POST /doctor/pregnancy/:id/end - End a pregnancy
func EndPregnancy(c *gin.Context) {
	id := c.Param("id")
	var payload struct {
		DeliveryDate   time.Time `json:"delivery_date"`
		DeliveryMethod string    `json:"delivery_method"`
		BirthWeight    float64   `json:"birth_weight"`
		Sex            string    `json:"sex"`
		DeliveryPlace  string    `json:"delivery_place"`
		Complications  string    `json:"complications"`
		ChildStatus    string    `json:"child_status"`
		GestationalAge int       `json:"gestational_age"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()

	// Find the pregnancy
	var pregnancy entity.Pregnancy
	if err := db.First(&pregnancy, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pregnancy not found"})
		return
	}

	if pregnancy.Status != "Active" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pregnancy is not active"})
		return
	}

	// Begin transaction
	tx := db.Begin()

	// 1. Update Pregnancy Status to Ended
	if err := tx.Model(&pregnancy).Update("status", "Ended").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update pregnancy status"})
		return
	}

	// 2. Create PreviousPregnancy record
	prevPregnancy := entity.PreviousPregnancy{
		PregnantWomanID: pregnancy.PregnantWomanID,
		PregnancyNo:     pregnancy.PregnancyNo,
		DeliveryDate:    payload.DeliveryDate,
		GestationalAge:  payload.GestationalAge,
		DeliveryMethod:  payload.DeliveryMethod,
		BirthWeight:     payload.BirthWeight,
		Sex:             payload.Sex,
		DeliveryPlace:   payload.DeliveryPlace,
		Complications:   payload.Complications,
		ChildStatus:     payload.ChildStatus,
	}

	if err := tx.Create(&prevPregnancy).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create previous pregnancy record"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Pregnancy ended successfully", "data": prevPregnancy})
}
