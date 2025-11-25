package entity

import "gorm.io/gorm"

type MedicalHistory struct {
	gorm.Model

	// FK -> PregnantWoman (P_ID)
	PregnantWomanID *uint          
	PregnantWoman   *PregnantWoman `gorm:"references:ID"`

	ChronicDiseases         string
	HeartDisease            bool
	Thyroid                 bool
	OtherDiseases           string
	SurgeryHistory          string // C-Section, Other
	OtherSurgery            string
	GeneticDiseases         string
	DrugAllergies           string

	// Family History
	FamilyHistoryHT          bool
	FamilyHistoryDiabetes    bool
	FamilyHistoryThalassemia bool
	FamilyHistoryCongenital  bool
	OtherFamilyHistory       string

	// Contraception
	ContraceptionBeforeMethod   string
	ContraceptionBeforeDuration string
	ContraceptionLastMethod     string
	ContraceptionLastDuration   string

	// Menstruation
	MenstrualCycle     int    // Every X days
	MenstrualDuration  int    // Lasts X days
	MenstrualCondition string // Regular/Irregular
}
