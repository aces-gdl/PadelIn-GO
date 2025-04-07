package server

import (
	"PadelIn/controllers"
	"PadelIn/middleware"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	Origins := strings.Split(os.Getenv("ORIGINS"), ",")
	corsConfig := cors.Config{
		AllowOrigins:     Origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.Use(gin.Recovery())

	// Security
	router.POST("/v1/security/signup", controllers.Signup)
	router.POST("/v1/security/login", controllers.Login)
	router.POST("/v1/security/logout", controllers.Logout)
	router.GET("/v1/security/validate", middleware.RequireAuth, controllers.Validate)
	router.POST("/v1/security/change-password", middleware.RequireAuth, controllers.ChangePassword)

	// Institution routes
	router.POST("/v1/institutions", middleware.RequireAuth, controllers.CreateInstitution)
	router.GET("/v1/institutions/:id", middleware.RequireAuth, controllers.GetInstitution)
	router.GET("/v1/institutions", middleware.RequireAuth, controllers.GetAllInstitutions)
	router.PUT("/v1/institutions/:id", middleware.RequireAuth, controllers.UpdateInstitution)
	router.DELETE("/v1/institutions/:id", middleware.RequireAuth, controllers.DeleteInstitution)

	// Payment routes
	router.POST("/v1/payments", middleware.RequireAuth, controllers.CreatePayment)
	router.GET("/v1/payments/:id", middleware.RequireAuth, controllers.GetPayment)
	router.GET("/v1/payments", middleware.RequireAuth, controllers.ListPayments)
	router.PUT("/v1/payments/:id", middleware.RequireAuth, controllers.UpdatePayment)
	router.DELETE("/v1/payments/:id", middleware.RequireAuth, controllers.DeletePayment)

	// clubs
	router.POST("/v1/catalogs/club", middleware.RequireAuth, controllers.PostClub)
	router.GET("/v1/catalogs/clubs", middleware.RequireAuth, controllers.GetClubs)

	// Tournaments
	router.GET("/v1/catalogs/tournaments", middleware.RequireAuth, controllers.GetTournaments)
	router.GET("/v1/catalogs/tournament", middleware.RequireAuth, controllers.GetTournament)
	router.POST("/v1/catalogs/tournaments", middleware.RequireAuth, controllers.PostTournaments)

	return router
}
