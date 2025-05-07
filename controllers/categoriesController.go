package controllers

import (
	"PadelIn/initializers"
	"PadelIn/middleware"
	"PadelIn/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostCategory(c *gin.Context) {

	var body models.TournamentCategory

	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al leer body...",
		})
		return
	}
	category := models.TournamentCategory{
		Description: body.Description,
		Level:       body.Level,
		Active:      body.Active,
		BgColor:     body.BgColor,
		TextColor:   body.TextColor,
	}
	category.InstitutionID = int(middleware.CurrentUser.InstitutionID)
	result := initializers.DB.Create(&category)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al crear usuario... ",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": category})
}

func GetCatgories(c *gin.Context) {

	CategoryID := c.DefaultQuery("CategoryID", "")

	var categories []models.TournamentCategory

	SQLStr := fmt.Sprintf(`Select * from tournament_categories where institution_id = %d `, middleware.CurrentUser.InstitutionID)
	if CategoryID != "" {
		SQLStr += fmt.Sprintf(` and id = %s`, CategoryID)
	}

	SQLStr += " order by level asc, description asc "
	results := initializers.DB.Debug().Raw(SQLStr).Scan(&categories)
	if results.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "results": len(categories), "data": categories})
}

func PutCategory(c *gin.Context) {
	// Get the category ID from the URL
	categoryID := c.Param("id")

	// Bind the request body to a struct
	var body struct {
		Description string `json:"description"`
		Level       int    `json:"level"`
		Active      bool   `json:"active"`
		BgColor     string `json:"bgColor"`
		TextColor   string `json:"textColor"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Find the category
	var category models.TournamentCategory
	result := initializers.DB.First(&category, categoryID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Check if the category belongs to the current user's institution
	if category.InstitutionID != int(middleware.CurrentUser.InstitutionID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this category"})
		return
	}

	// Update the category
	category.Description = body.Description
	category.Level = body.Level
	category.Active = body.Active
	category.BgColor = body.BgColor
	category.TextColor = body.TextColor
	result = initializers.DB.Save(&category)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": category})
}

func DeleteCategory(c *gin.Context) {
	// Get the category ID from the URL
	categoryID := c.Param("id")

	// Find the category
	var category models.TournamentCategory
	result := initializers.DB.First(&category, categoryID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Check if the category belongs to the current user's institution
	if category.InstitutionID != int(middleware.CurrentUser.InstitutionID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this category"})
		return
	}

	// Delete the category
	result = initializers.DB.Delete(&category)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Category deleted successfully"})
}
