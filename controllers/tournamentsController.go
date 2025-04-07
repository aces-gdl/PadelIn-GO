package controllers

import (
	"PadelIn/initializers"
	"PadelIn/middleware"
	"PadelIn/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetTournaments(c *gin.Context) {

	fmt.Println("User ", middleware.CurrentUser.ID)

	var tournaments []models.Tournament
	sqlStatement := `Select * from tournaments where id_institution =%d order by id desc`
	sqlStatement = fmt.Sprintf(sqlStatement, middleware.CurrentUser.InstitutionID)
	results := initializers.DB.Debug().Raw(sqlStatement).Scan(&tournaments)
	if results.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "results": len(tournaments), "data": tournaments})
}

func GetTournament(c *gin.Context) {
	var TournamentID = c.DefaultQuery("TournamentID", "")

	if TournamentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "TournamentID is required..."})
		return
	}
	var tournament models.Tournament
	var sqlStatement = `Select * from tournaments`
	if TournamentID != "" {
		sqlStatement += ` where id = ` + TournamentID
	}
	sqlStatement += ` order by id desc`
	results := initializers.DB.Debug().Raw(sqlStatement).Scan(&tournament)

	if results.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "TournamentID Not Found..."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "results": 1, "data": tournament})
}

func ParseTime(timeOnlySrt string) time.Time {
	timeArray := strings.Split(timeOnlySrt, ":")
	hour, _ := strconv.Atoi(timeArray[0])
	minute, _ := strconv.Atoi(timeArray[1])
	myLocation, _ := time.LoadLocation("America/Mexico_City")
	result := time.Date(2023, 10, 02, hour, minute, 0, 0, myLocation)

	return result
}

func PostTournaments(c *gin.Context) {

	var body models.Tournament
	resultTest := c.Bind(&body)
	if resultTest != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al leer body...",
		})
		return
	}

	result := initializers.DB.Create(&body)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al crear torneo... ",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "results": 1, "data": body})
}
