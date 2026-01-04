package entity

import (
	"time"

	"gorm.io/gorm"
)

type Vaccination struct {
	gorm.Model

	// FK -> PregnantWoman (P_ID)
	PregnantWomanID *uint          `gorm:"column:p_id" json:"PregnantWomanID" valid:"required~กรุณาเลือกมารดา"`
	PregnantWoman   *PregnantWoman `gorm:"references:ID" valid:"-"`

	IsPreviouslyVaccinated bool      `json:"IsPreviouslyVaccinated"`
	PreviousDoses          int       `json:"PreviousDoses"`
	LastPreviousDateYear   *time.Time `json:"LastPreviousDateYear"`
	Dose1DateDuringPreg    *time.Time `json:"Dose1DateDuringPreg"`
	Dose2DateDuringPreg    *time.Time `json:"Dose2DateDuringPreg"`
	Dose3DateDuringPreg    *time.Time `json:"Dose3DateDuringPreg"`
	IsHistoryUnknown       bool      `json:"IsHistoryUnknown"`
	ReasonForNotVaccinating string    `json:"ReasonForNotVaccinating"`
	Remarks                string    `json:"Remarks"`

	// FK -> VaccineType
	VaccineTypeID *uint        `json:"VaccineTypeID" valid:"-"`
	VaccineType   *VaccineType `gorm:"references:ID" valid:"-"`

	// FK -> VacDose (VD_ID)
	VacDoseID *uint    `gorm:"column:vd_id" valid:"-"`
	VacDose   *VacDose `gorm:"references:ID" valid:"-"`

	Doses []VaccineDose `gorm:"foreignKey:VaccinationID"`
}
