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

	// Category
	router.POST("/v1/catalogs/categories", middleware.RequireAuth, controllers.PostCategory)
	router.GET("/v1/catalogs/categories", middleware.RequireAuth, controllers.GetCatgories)
	router.PUT("/v1/catalogs/categories/:id", middleware.RequireAuth, controllers.PutCategory)
	router.DELETE("/v1/catalogs/categories/:id", middleware.RequireAuth, controllers.DeleteCategory)

	// Players By Tournament routes
	router.POST("/v1/tournaments/register", middleware.RequireAuth, controllers.PostRegisterUserForTournament)
	router.DELETE("/v1/tournaments/unregister", middleware.RequireAuth, controllers.DeleteRegisterUserForTournament)
	router.GET("/v1/tournaments/players", middleware.RequireAuth, controllers.GetRegisteredUsersByTournament)
	router.PUT("/v1/tournaments/players/:id", middleware.RequireAuth, controllers.UpdatePlayerByTournament)

	// Teams
	router.GET("/v1/tournament/isplayertaken", middleware.RequireAuth, controllers.GetIsPlayerTaken)
	router.POST("/v1/tournament/enrolledteams", middleware.RequireAuth, controllers.PostEnrolledTeams)
	router.PUT("/v1/tournament/playerbytournament/:id", controllers.UpdatePlayerByTournament)
	router.GET("/v1/tournament/playerbytournament", middleware.RequireAuth, controllers.GetEnrolledTeams)
	router.POST("/v1/tournament/createblankteam", middleware.RequireAuth, controllers.PostBlankTeam)
	router.DELETE("/v1/tournament/deleteteam/:id", middleware.RequireAuth, controllers.DeleteEnrolledTeam)

	// Stripe
	router.GET("/v1/stripe/config", controllers.StripeConfig)
	router.GET("/v1/stripe/get-payment-amount", controllers.GetPaymentAmount)
	router.POST("/v1/stripe/create-payment-intent", controllers.CreatePaymentIntent)

	// Player By Category
	router.GET("/v1/catalogs/playersbycategory", middleware.RequireAuth, controllers.PlayersByCategory)

	return router
}
