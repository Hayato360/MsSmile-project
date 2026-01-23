package entity

import (
	"time"

	"gorm.io/gorm"
)

type LabResult struct {
	gorm.Model

	// FK -> Pregnancy
	PregnancyID *uint      `valid:"required~กรุณาเลือกครรภ์"`
	Pregnancy   *Pregnancy `gorm:"references:ID" valid:"-"`

	TestDate     time.Time    `json:"TestDate"`
	Hct          float64      `json:"Hct"` // Hematocrit
	Hb           float64      `json:"Hb"`  // Hemoglobin
	HbTyping     string       `json:"HbTyping"`
	OtherRemarks string       `json:"OtherRemarks"`
	FilePath     string       `json:"FilePath"` // Path to the uploaded PDF file

	// FK -> ผลตรวจ
	DCPResultID     *uint        `valid:"-" json:"DCPResultID"`
	DCPResult       *CheckResult `gorm:"foreignKey:DCPResultID" valid:"-" json:"DCPResult"`
	AntiHIVResultID *uint        `valid:"-" json:"AntiHIVResultID"`
	AntiHIVResult   *CheckResult `gorm:"foreignKey:AntiHIVResultID" valid:"-" json:"AntiHIVResult"`
}
