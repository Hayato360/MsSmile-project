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
		&entity.Pregnancy{},
		&entity.PreviousPregnancy{},
		&entity.MedicalHistory{},
		&entity.Vaccination{},
		&entity.AntenatalVisit{},
		&entity.LabResult{},
		&entity.FetalKickCount{},
		&entity.Appointment{},
		&entity.VaccineDose{},
	)

	hashedPassword, _ := HashPassword("123456")

	// Create Doctor
	Doctor := entity.Doctor{
		Username:        "Doctor",
		Password:        hashedPassword,
		Email:           "Doctor@gmail.com",
		FullName:        "Doctor D",
		PhoneNumber:     "0655765587",
		CitizenID:       "1111111111111",
		HospitalCode:    "HOS001",
		DoctorLicenseNo: "DOC00001",
		IsRegistered:    true, // Already registered
	}
	db.FirstOrCreate(&Doctor, &entity.Doctor{Username: "Doctor"})

	// Seed Pre-verified Doctors (NOT registered yet)
	// These doctors can register using their credentials
	PreVerifiedDoc1 := entity.Doctor{
		CitizenID:       "1234567890123",
		HospitalCode:    "HOS001",
		DoctorLicenseNo: "DOC12345",
		FullName:        "นพ.สมชาย ใจดี",
		Email:           "somchai@hospital.th",
		PhoneNumber:     "0812345678",
		IsRegistered:    false,
		// Username and Password are empty (will be set during registration)
	}
	db.FirstOrCreate(&PreVerifiedDoc1, &entity.Doctor{
		CitizenID:       "1234567890123",
		HospitalCode:    "HOS001",
		DoctorLicenseNo: "DOC12345",
	})

	PreVerifiedDoc2 := entity.Doctor{
		CitizenID:       "9876543210987",
		HospitalCode:    "HOS001",
		DoctorLicenseNo: "DOC67890",
		FullName:        "พญ.สมหญิง รักษา",
		Email:           "somying@hospital.th",
		PhoneNumber:     "0898765432",
		IsRegistered:    false,
	}
	db.FirstOrCreate(&PreVerifiedDoc2, &entity.Doctor{
		CitizenID:       "9876543210987",
		HospitalCode:    "HOS001",
		DoctorLicenseNo: "DOC67890",
	})

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
