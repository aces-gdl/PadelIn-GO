package controllers

import (
	"PadelIn/initializers"
	"PadelIn/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateInstitution handles the creation of a new institution
func CreateInstitution(c *gin.Context) {
	var body struct {
		Name        string
		Description string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	institution := models.Institution{
		Name:        body.Name,
		Description: body.Description,
	}

	result := initializers.DB.Create(&institution)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create institution",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"institution": institution,
	})
}

// GetInstitution retrieves a specific institution by ID
func GetInstitution(c *gin.Context) {
	id := c.Param("id")

	var institution models.Institution
	result := initializers.DB.First(&institution, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Institution not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"institution": institution,
	})
}

// GetAllInstitutions retrieves all institutions
func GetAllInstitutions(c *gin.Context) {
	var institutions []models.Institution
	result := initializers.DB.Find(&institutions)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve institutions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"institutions": institutions,
	})
}

// UpdateInstitution updates an existing institution
func UpdateInstitution(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		Name        string
		Description string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var institution models.Institution
	result := initializers.DB.First(&institution, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Institution not found",
		})
		return
	}

	initializers.DB.Model(&institution).Updates(models.Institution{
		Name:        body.Name,
		Description: body.Description,
	})

	c.JSON(http.StatusOK, gin.H{
		"institution": institution,
	})
}

// DeleteInstitution deletes an institution
func DeleteInstitution(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	result := initializers.DB.Delete(&models.Institution{}, id)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete institution",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Institution not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Institution deleted successfully",
	})
}
