package entity

import "gorm.io/gorm"

type Doctor struct {
	gorm.Model
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	
	// **ฟิลด์ที่เพิ่มเข้ามาสำหรับการล็อกอิน**
	Username 	string 	`gorm:"uniqueIndex" json:"username"`
	Password 	string `json:"-"`

	VisitDoctors []VisitDoctor `gorm:"foreignKey:DoctorID"`
}