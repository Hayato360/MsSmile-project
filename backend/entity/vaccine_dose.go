package entity

import (
	"time"

	"gorm.io/gorm"
)

type VaccineDose struct {
	gorm.Model
	VaccinationID uint
	Vaccination   *Vaccination `gorm:"foreignKey:VaccinationID"`

	DoseDate time.Time
	DoseNo   int
}
