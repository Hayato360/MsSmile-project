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
		protected.GET("/previous-pregnancies/pregnant-woman/:id", controller.GetPreviousPregnanciesByPregnantWomanID)

		// Antenatal Visit Routes
		protected.GET("/antenatal-visits/pregnancy/:id", controller.GetAntenatalVisitsByPregnancyID)
		protected.POST("/antenatal-visits", controller.CreateAntenatalVisit)

		// Doctor Routes
		protected.GET("/doctor/patients", controller.GetDoctorPatients)
		protected.GET("/doctor/patients/:id", controller.GetDoctorPatientDetail)
		protected.GET("/doctor/patient/:patientId/visits", controller.GetPatientVisits)
		protected.POST("/doctor/antenatal-visit", controller.DoctorCreateAntenatalVisit)
		protected.POST("/doctor/pregnancy", controller.DoctorCreatePregnancy)
		protected.POST("/doctor/pregnancy/:id/end", controller.EndPregnancy)
		protected.POST("/doctor/patient/:id/appointment", controller.DoctorCreateAppointment)

		// Doctor Health Data Routes
		protected.POST("/doctor/medical-history", controller.DoctorCreateMedicalHistory)
		protected.GET("/doctor/patient/:patientId/medical-history", controller.GetPatientMedicalHistory)
		protected.POST("/doctor/lab-result", controller.DoctorCreateLabResult)
		protected.POST("/doctor/vaccination", controller.DoctorCreateVaccination)
		protected.GET("/vaccine-types", controller.ListVaccineTypes)
		
		protected.POST("/doctor/previous-pregnancy", controller.DoctorCreatePreviousPregnancy)
		protected.GET("/doctor/patient/:patientId/previous-pregnancies", controller.GetPreviousPregnancies)
		protected.GET("/doctor/pregnancy/:pregnancyId/lab-results", controller.GetLabResultsByPregnancyID)
		
		// Static file serving for uploads
		r.Static("/uploads", "./uploads")

		// Profile Routes
		protected.PUT("/profile/husband", controller.UpdateHusband)
		protected.PUT("/profile/personal", controller.UpdatePersonalProfile)
		protected.PUT("/profile/medical-history", controller.UpdateMyMedicalHistory)
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


