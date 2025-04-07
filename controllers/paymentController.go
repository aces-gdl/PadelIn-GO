package controllers

import (
	"PadelIn/initializers"
	"PadelIn/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PaymentActivity struct {
	Activity []struct {
		TransactionDate string `json:"transactionDate"`
		Notes           string `json:"notes"`
	} `json:"activity"`
}

type Activity = struct {
	TransactionDate string `json:"transactionDate"`
	Notes           string `json:"notes"`
}

// CreatePayment handles the creation of a new payment
func CreatePayment(c *gin.Context) {

	var body models.Payment

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	myActivity := `{
	    "Activity": [
			{ 
			"transactionDate": "2022-01-01T15:04:05Z", 
			"notes": "Transaccion Pendiente" 
			}
		]
	}`

	var activity PaymentActivity
	err := json.Unmarshal([]byte(myActivity), &activity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse activity data"})
		return
	}
	jsonData, _ := json.Marshal(activity)
	data := string(jsonData)

	body.MetaData = data

	result := initializers.DB.Create(&body)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"payment": body})
}

// GetPayment retrieves a specific payment by ID
func GetPayment(c *gin.Context) {
	id := c.Param("id")

	var payment models.Payment
	result := initializers.DB.First(&payment, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment": payment})
}

// ListPayments retrieves a list of payments with optional filtering
func ListPayments(c *gin.Context) {
	var payments []models.Payment
	query := initializers.DB.Model(&models.Payment{})

	// Add filters based on query parameters
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	// Add more filters as needed

	result := query.Find(&payments)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve payments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payments": payments})
}

// UpdatePayment updates an existing payment
func UpdatePayment(c *gin.Context) {
	id := c.Param("id")

	var body models.Payment

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var payment models.Payment
	result := initializers.DB.First(&payment, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}
	// TODO: En el body deberia de venir que tipo de actualizacion es : paid, cancelled, refunded, pending, processing

	// Update fields
	if body.Status == "cancelled" && payment.Status != "paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se puede cancelar un pago que haya sido pagado"})
		return
	}
	payment.Status = body.Status

	if body.Status == "paid" {
		now := time.Now()
		payment.PaymentDate = &now

		var activity PaymentActivity
		err := json.Unmarshal([]byte(payment.MetaData), &activity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse activity data"})
			return
		}
		activity.Activity = append(activity.Activity, Activity{TransactionDate: now.Format(time.RFC3339), Notes: "Pago realizado"})

		jsonData, _ := json.Marshal(activity)

		data := string(jsonData)
		payment.PaymentMethod = "card"
		payment.PaymentVendorID = body.PaymentVendorID
		payment.MetaData = data
	}

	if body.Status == "cancelled" {
		var activity PaymentActivity
		err := json.Unmarshal([]byte(payment.MetaData), &activity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse activity data"})
			return
		}
		activity.Activity = append(activity.Activity, Activity{TransactionDate: time.Now().Format(time.RFC3339), Notes: "Pago cancelado"})

		jsonData, _ := json.Marshal(activity)

		data := string(jsonData)
		payment.MetaData = data
	}
	payment.RefundStatus = body.RefundStatus
	payment.RefundDate = body.RefundDate

	result = initializers.DB.Save(&payment)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment": payment})
}

// DeletePayment soft deletes a payment
func DeletePayment(c *gin.Context) {
	id := c.Param("id")

	result := initializers.DB.Delete(&models.Payment{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment deleted successfully"})
}
