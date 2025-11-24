package controller

import (
	"net/http"
	"time"

	"github.com/bestiesmile1845/Projecteiei/config"
	"github.com/bestiesmile1845/Projecteiei/entity"
	"github.com/gin-gonic/gin"
)

// POST /kick-counts
func CreateFetalKickCount(c *gin.Context) {
	var kickCount entity.FetalKickCount

	// Bind JSON payload
	if err := c.ShouldBindJSON(&kickCount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()

	// Validate PregnancyID
	if kickCount.PregnancyID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PregnancyID is required"})
		return
	}

	// Define start and end of the day for the incoming record (in UTC to match DB storage usually, or just use the date part)
	// Assuming CountDate comes in as a valid time.Time
	// We want to match any record that falls on the same calendar day.
	
	// Truncate to 24-hour window based on the input date
	// Note: This assumes the input CountDate is set to the correct day.
	// Ideally, we should handle timezones, but for simplicity in this project, 
	// we'll construct a range covering the whole day of the given CountDate.
	
	year, month, day := kickCount.CountDate.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, kickCount.CountDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Check if a record exists for this pregnancy on this date
	var existingRecord entity.FetalKickCount
	err := db.Where("pregnancy_id = ? AND count_date >= ? AND count_date < ?", kickCount.PregnancyID, startOfDay, endOfDay).First(&existingRecord).Error

	if err == nil {
		// Record exists, update it
		existingRecord.KickCountMorning = kickCount.KickCountMorning
		existingRecord.KickCountLunch = kickCount.KickCountLunch
		existingRecord.KickCountEvening = kickCount.KickCountEvening
		// Update CountDate to the new one if needed, or keep the original
		// existingRecord.CountDate = kickCount.CountDate 

		if err := db.Save(&existingRecord).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Updated success", "data": existingRecord})
	} else {
		// Record does not exist, create new
		if err := db.Create(&kickCount).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Created success", "data": kickCount})
	}
}

// GET /kick-counts/pregnancy/:id
func GetFetalKickCountsByPregnancyID(c *gin.Context) {
	id := c.Param("id")
	var kickCounts []entity.FetalKickCount

	db := config.DB()

	// Find records by PregnancyID and Sort by Date
	if err := db.Where("pregnancy_id = ?", id).Order("count_date ASC").Find(&kickCounts).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kickCounts)
}
