package config

import (
	"fmt"
	"os"
	"time"

	"github.com/bestiesmile1845/Projecteiei/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func DB() *gorm.DB {
	return db
}

func ConnectionDB() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "Mother.db"
	}
	database, err := gorm.Open(sqlite.Open(dbPath+"?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("connected database")
	db = database
}

func SetupDatabase() {
	// AutoMigrate will create the tables and columns if they do not exist
	db.AutoMigrate(
		&entity.PregnantWoman{},
		&entity.Doctor{},
		&entity.PreviousPregnancy{},
		&entity.MedicalHistory{},
		&entity.Vaccination{},
		&entity.AntenatalVisit{},
		&entity.LabResult{},
		&entity.FetalKickCount{},
		&entity.Appointment{},
	)

	hashedPassword, _ := HashPassword("123456")

	// Create Doctor
	Doctor := entity.Doctor{
		Username:    "Doctor",
		Password:    hashedPassword,
		Email:       "Doctor@gmail.com",
		FullName:    "Doctor D",
		PhoneNumber: "0655765587",
	}
	db.FirstOrCreate(&Doctor, &entity.Doctor{Username: "Doctor"})

	// Create Pregnant Woman (Not Pregnant yet)
	Woman := entity.PregnantWoman{
		Username:    "Mommy",
		Password:    hashedPassword,
		Email:       "Mommy@gmail.com",
		FullName:    "Mommy M",
		PhoneNumber: "0812345678",
		Age:         25,
	}
	db.FirstOrCreate(&Woman, &entity.PregnantWoman{Username: "Mommy"})

	// Seed Vaccine Types
	vaccineTypes := []entity.VaccineType{
		{Name: "บาดทะยัก-คอตีบ (dT)"},
		{Name: "ไข้หวัดใหญ่ (Influenza)"},
		{Name: "โควิด 19 (Covid-19)"},
	}
	for _, vt := range vaccineTypes {
		db.FirstOrCreate(&vt, entity.VaccineType{Name: vt.Name})
	}

	// Create Appointment and Assign to Mommy
	// 25 Nov 2025 09:00:00
	apptDate, _ := time.Parse("2006-01-02 15:04:05", "2025-11-25 09:00:00")
	Appt := entity.Appointment{
		AppointmentDate: apptDate,
		Title:           "นัดตรวจครรภ์ครั้งถัดไป",
		Location:        "อาคารผู้ป่วยนอก",
	}
	db.FirstOrCreate(&Appt, &entity.Appointment{Title: "นัดตรวจครรภ์ครั้งถัดไป"})

	// Update Mommy
	var mommy entity.PregnantWoman
	if err := db.Where("username = ?", "Mommy").First(&mommy).Error; err == nil {
		mommy.AppointmentID = &Appt.ID
		db.Save(&mommy)
	}
}
