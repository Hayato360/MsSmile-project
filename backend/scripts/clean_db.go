package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	cwd, _ := os.Getwd()
	fmt.Printf("Current Working Directory: %s\n", cwd)

	// Check if Mother.db exists in current directory
	dbPath := "Mother.db"
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Printf("Mother.db not found in %s\n", cwd)
		// Try parent directory
		dbPath = "../Mother.db"
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			fmt.Println("Error: Mother.db not found in current or parent directory.")
			return
		}
	}

	fmt.Printf("Found database at: %s\n", dbPath)
	
	// Use absolute path to be safe
	absPath, _ := filepath.Abs(dbPath)
	fmt.Printf("Absolute path: %s\n", absPath)

	db, err := gorm.Open(sqlite.Open(absPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// List of tables to clean (Order matters for foreign keys if enforced, but we'll disable them temporarily)
	tables := []string{
		"fetal_kick_counts",
		"lab_results",
		"antenatal_visits",
		"vaccinations",
		"medical_histories",
		"previous_pregnancies",
		"pregnancies",
		"pregnant_women",
	}

	// Disable Foreign Keys
	db.Exec("PRAGMA foreign_keys = OFF")

	for _, table := range tables {
		err := db.Exec("DELETE FROM " + table).Error
		if err != nil {
			fmt.Printf("Error cleaning table %s: %v\n", table, err)
		} else {
			fmt.Printf("Cleaned table: %s\n", table)
			// Reset Auto Increment
			db.Exec("DELETE FROM sqlite_sequence WHERE name = ?", table)
		}
	}

	// Enable Foreign Keys
	db.Exec("PRAGMA foreign_keys = ON")

	fmt.Println("---------------------------------------------------")
	fmt.Println("Database cleanup completed successfully.")
	fmt.Println("All data removed except 'doctors' table.")
	fmt.Println("---------------------------------------------------")
}
