package entity

import (
	"time"

	"gorm.io/gorm"
)

type PreviousPregnancy struct {
	gorm.Model

	// FK -> PregnantWoman
	PregnantWomanID *uint
	PregnantWoman   *PregnantWoman `gorm:"references:ID"`

	PregnancyNo    int       `json:"pregnancy_no"`
	DeliveryDate   time.Time `json:"delivery_date"`
	GestationalAge int       `json:"gestational_age"` // Weeks
	DeliveryMethod string    `json:"delivery_method"` // Normal, C-Section, Vacuum, etc.
	BirthWeight    float64   `json:"birth_weight"`    // kg
	Sex            string    `json:"sex"`             // Male, Female
	DeliveryPlace  string    `json:"delivery_place"`
	Complications  string    `json:"complications"`
	ChildStatus    string    `json:"child_status"` // Healthy, Deceased, etc.
}
