package controller

import (
	"net/http"
	"fmt"

	"github.com/bestiesmile1845/Projecteiei/config"
	"github.com/bestiesmile1845/Projecteiei/entity"
	"github.com/gin-gonic/gin"
)

// GET /vaccinations/pregnant-woman/:id
func GetVaccinationsByPregnantWomanID(c *gin.Context) {
	id := c.Param("id")
	var vaccinations []entity.Vaccination

	db := config.DB()

	if err := db.Preload("VaccineType").Preload("VacDose").Preload("Doses").Where("p_id = ?", id).Find(&vaccinations).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vaccinations)
}

// PUT /vaccinations/:id
func UpdateVaccination(c *gin.Context) {
	id := c.Param("id")
	var vaccination entity.Vaccination

	db := config.DB()

	if err := db.First(&vaccination, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vaccination not found"})
		return
	}

	if err := c.ShouldBindJSON(&vaccination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&vaccination).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated success", "data": vaccination})
}

// GET /vaccine-types
func ListVaccineTypes(c *gin.Context) {
	var vaccineTypes []entity.VaccineType
	db := config.DB()
	if err := db.Find(&vaccineTypes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": vaccineTypes})
}

// POST /doctor/vaccination
func DoctorCreateVaccination(c *gin.Context) {
	var payload entity.Vaccination
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

    // DEBUG LOGGING
    fmt.Printf("DEBUG POST VACCINATION: Payload: %+v\n", payload)
    if payload.PregnantWomanID != nil {
        fmt.Printf("DEBUG: PregnantWomanID: %d\n", *payload.PregnantWomanID)
    } else {
        fmt.Println("DEBUG: PregnantWomanID is NIL")
    }
    if payload.VaccineTypeID != nil {
        fmt.Printf("DEBUG: VaccineTypeID: %d\n", *payload.VaccineTypeID)
    } else {
        fmt.Println("DEBUG: VaccineTypeID is NIL")
    }

	db := config.DB()
	var existingVaccination entity.Vaccination

	// Check if exists
	if err := db.Where("p_id = ? AND vaccine_type_id = ?", payload.PregnantWomanID, payload.VaccineTypeID).First(&existingVaccination).Error; err == nil {
		// Exists -> Update
		existingVaccination.IsPreviouslyVaccinated = payload.IsPreviouslyVaccinated
		existingVaccination.PreviousDoses = payload.PreviousDoses
		existingVaccination.LastPreviousDateYear = payload.LastPreviousDateYear
		existingVaccination.Dose1DateDuringPreg = payload.Dose1DateDuringPreg
		existingVaccination.Dose2DateDuringPreg = payload.Dose2DateDuringPreg
		existingVaccination.Dose3DateDuringPreg = payload.Dose3DateDuringPreg
		existingVaccination.IsHistoryUnknown = payload.IsHistoryUnknown
		existingVaccination.ReasonForNotVaccinating = payload.ReasonForNotVaccinating
		existingVaccination.Remarks = payload.Remarks

		// Replace Doses
		if err := db.Model(&existingVaccination).Association("Doses").Replace(payload.Doses); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := db.Save(&existingVaccination).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		// Reload to get updated data including doses
		db.Preload("Doses").First(&existingVaccination, existingVaccination.ID)
		
		c.JSON(http.StatusOK, gin.H{"message": "Updated success", "data": existingVaccination})
		return
	}

	// Not exists -> Create
	if err := db.Create(&payload).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Created success", "data": payload})
}

// DELETE /vaccinations/:id
func DeleteVaccination(c *gin.Context) {
	id := c.Param("id")
	db := config.DB()

	var vaccination entity.Vaccination
	if err := db.Preload("Doses").First(&vaccination, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vaccination record not found"})
		return
	}

	// Delete association (Doses)
	// GORM cascading delete might work if configured, but let's be explicit or rely on Unscoped/Association clear
	if err := db.Model(&vaccination).Association("Doses").Clear(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete doses assoc"})
		return
	}

	// Also delete the Doses records themselves if they are orphan? 
	// The Association Clear just removes FK. We might want to remove them.
	// Actually, if we delete Vaccination, GORM might delete children if Constraint is set.
	// But let's just delete the vaccination.
	
	if err := db.Select("Doses").Delete(&vaccination).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted success"})
}
