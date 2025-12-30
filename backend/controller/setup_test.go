package controller_test

import (
	"os"
	"github.com/bestiesmile1845/Projecteiei/config"
	"github.com/gin-gonic/gin"
)

func setupTestDB() {
	// Set DB_PATH to :memory: to use an in-memory database
	os.Setenv("DB_PATH", ":memory:")
	config.ConnectionDB()
	config.SetupDatabase() // Migrates and seeds data
	gin.SetMode(gin.TestMode)
}
