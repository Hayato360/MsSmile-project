package main

import (
	"github.com/bestiesmile1845/Projecteiei/config"
	"github.com/bestiesmile1845/Projecteiei/controller"
	"github.com/bestiesmile1845/Projecteiei/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {

	// Initialize database connection
	config.ConnectionDB()
	config.SetupDatabase()

	r := gin.Default()

	r.Use(CORSMiddleware())

	r.POST("/login", controller.Login)
	r.POST("/register", controller.CreatePregnantWoman)

	// Protected Routes
	protected := r.Group("/")
	protected.Use(middlewares.Authorizes())
	{
		protected.GET("/me", controller.GetMe)

		// Kick Count Routes
		protected.POST("/kick-counts", controller.CreateFetalKickCount)
		protected.GET("/kick-counts/pregnancy/:id", controller.GetFetalKickCountsByPregnancyID)

		// Vaccination Routes
		protected.GET("/vaccinations/pregnant-woman/:id", controller.GetVaccinationsByPregnantWomanID)
		protected.PUT("/vaccinations/:id", controller.UpdateVaccination)

		// Medical History Routes
		protected.GET("/medical-histories/pregnant-woman/:id", controller.GetMedicalHistoryByPregnantWomanID)
		protected.PUT("/medical-histories/:id", controller.UpdateMedicalHistory)

		// Antenatal Visit Routes
		protected.GET("/antenatal-visits/pregnancy/:id", controller.GetAntenatalVisitsByPregnancyID)
		protected.POST("/antenatal-visits", controller.CreateAntenatalVisit)

		// Doctor Routes
		protected.GET("/doctor/patients", controller.GetDoctorPatients)
		protected.GET("/doctor/patients/:id", controller.GetDoctorPatientDetail)
		protected.GET("/doctor/patient/:patientId/visits", controller.GetPatientVisits)
		protected.POST("/doctor/antenatal-visit", controller.DoctorCreateAntenatalVisit)
		protected.POST("/doctor/pregnancy", controller.DoctorCreatePregnancy)

		// Doctor Health Data Routes
		protected.POST("/doctor/medical-history", controller.DoctorCreateMedicalHistory)
		protected.GET("/doctor/patient/:patientId/medical-history", controller.GetPatientMedicalHistory)
		protected.POST("/doctor/lab-result", controller.DoctorCreateLabResult)
		protected.POST("/doctor/vaccination", controller.DoctorCreateVaccination)

		// Profile Routes
		protected.PUT("/profile/husband", controller.UpdateHusband)
	}


	r.Run(":8081")

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Username")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}


