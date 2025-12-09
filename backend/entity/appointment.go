package entity

import (
	"time"

	"gorm.io/gorm"
)

type Appointment struct {
	gorm.Model
	AppointmentDate time.Time `json:"appointment_date"`
	Title           string    `json:"title"`
	Location        string    `json:"location"`

	VisitDoctorID *uint
	VisitDoctor   *VisitDoctor `gorm:"references:ID"`

	PregnantWomen []PregnantWoman `gorm:"foreignKey:AppointmentID"`
}
