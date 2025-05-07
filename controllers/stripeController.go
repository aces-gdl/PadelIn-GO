package controllers

import (
	"PadelIn/initializers"
	"PadelIn/middleware"
	"PadelIn/models"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/paymentintent"
)

func CreatePaymentIntent(c *gin.Context) {
	type BodyDefinition struct {
		TournamentID uint
		CategoryID   uint
	}

	var body BodyDefinition
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tournament models.Tournament
	result := initializers.DB.First(&tournament, body.TournamentID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Torneo no encontrado"})
		return
	}

	var category models.TournamentCategory
	result = initializers.DB.First(&category, body.CategoryID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Categoría no encontrada"})
		return
	}
	var amount int64 = int64(math.Round(tournament.InscriptionFee * 100))

	pk := os.Getenv("STRIPE_SECRET_KEY")
	stripe.Key = pk
	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(amount),
		Currency:           stripe.String("mxn"),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Description:        stripe.String("torneo golang"),
	}
	pi, err := paymentintent.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var payment models.Payment

	payment.InstitutionID = middleware.CurrentUser.InstitutionID
	payment.UserID = middleware.CurrentUser.ID
	payment.TournamentID = body.TournamentID
	payment.CategoryID = body.CategoryID
	payment.Amount = tournament.InscriptionFee
	payment.Currency = "MXN"
	payment.PaymentType = "torneo"
	payment.PaymentMethod = "card"
	payment.PaymentVendorID = "stripe"
	payment.TransactionID = pi.ID
	payment.Description = "Inscripcion a torneo :" + tournament.Description + " " + category.Description
	payment.CreatedAt = time.Now()
	payment.Status = "pending"
	initializers.DB.Create(&payment)

	c.JSON(http.StatusOK, gin.H{"clientSecret": pi.ClientSecret, "payment": payment})
}

func StripeConfig(c *gin.Context) {
	sk := os.Getenv("STRIPE_PUBLISHABLE_KEY")
	c.JSON(http.StatusOK, gin.H{"secret": sk})
}

func GetPaymentAmount(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"amount": 1099})
}

func SuccessfulPayment(c *gin.Context) {
	var body struct {
		PaymentID     string `json:"payment_id"`
		ClientSecret  string `json:"client_secret"`
		UserID        uint   `json:"user_id"`
		InstitutionID uint   `json:"institution_id"`
		TournamentID  uint   `json:"tournament_id"`
		CategoryID    uint   `json:"category_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var payment models.Payment

	result := initializers.DB.First(&payment, body.PaymentID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Pago no encontrado"})
		return
	}

	payment.Status = "paid"

	c.JSON(http.StatusOK, gin.H{"message": "Pago exitoso"})
}

func FailedPayment(c *gin.Context) {
	// Se usara para manejar el caso en el que el pago falle por:
	// Fondos insuficiente
	// Tarjeta Expirada
	// Tarjeta cancelada

	c.JSON(http.StatusOK, gin.H{"message": "Pago fallido"})
}
