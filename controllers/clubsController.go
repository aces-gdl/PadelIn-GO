package controllers

import (
	"PadelIn/initializers"
	"PadelIn/middleware"
	"PadelIn/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostClub(c *gin.Context) {
	var body models.Club

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al leer body...",
		})
		return
	}
	club := body
	club.InstitutionID = int(middleware.CurrentUser.InstitutionID)
	result := initializers.DB.Create(&club)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al crear usuario... ",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "results": 1, "data": club})
}

func GetClubs(c *gin.Context) {

	var clubs []models.Club
	results := initializers.DB.Debug().Find(&clubs, "institution_id =?", middleware.CurrentUser.InstitutionID)
	if results.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "results": len(clubs), "data": clubs})
}
