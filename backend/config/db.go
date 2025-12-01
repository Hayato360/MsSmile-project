package config

import (
	"fmt"
	"os"

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
}
