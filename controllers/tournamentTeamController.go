package controllers

import (
	"PadelIn/initializers"
	"PadelIn/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetIsPlayerTaken(c *gin.Context) {

	TournamentIDstr := c.DefaultQuery("TournamentID", "0")
	CategoryIDstr := c.DefaultQuery("CategoryID", "0")
	UserIDstr := c.DefaultQuery("UserID", "0")

	var team models.TournamentTeam
	whereStatement := `tournament_id=` + TournamentIDstr + ` and category_id = ` + CategoryIDstr + ` and ( member1_id = ` + UserIDstr + ` or member2_id = ` + UserIDstr + ` )`
	initializers.DB.Debug().Find(&team, whereStatement)
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": team})

}

func PostEnrolledTeams(c *gin.Context) {
	var body models.TournamentTeam

	bindResult := c.Bind(&body)
	if bindResult != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al leer body...",
		})
		return
	}

	tournamentTeam := models.TournamentTeam{
		Name:         body.Name,
		PlayLevel:    body.PlayLevel,
		Player1ID:    body.Player1ID,
		Name1:        body.Name1,
		Ranking1:     body.Ranking1,
		Player2ID:    body.Player2ID,
		Name2:        body.Name2,
		Ranking2:     body.Ranking2,
		CategoryID:   body.CategoryID,
		TournamentID: body.TournamentID,
	}

	result := initializers.DB.Create(&tournamentTeam)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al crear usuario... ",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ID": tournamentTeam.ID})
}

func PostBlankTeam(c *gin.Context) {
	var body models.TournamentTeam

	bindResult := c.Bind(&body)
	if bindResult != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al leer body...",
		})
		return
	}

	tournamentTeam := models.TournamentTeam{
		Name:         body.Name,
		CategoryID:   body.CategoryID,
		TournamentID: body.TournamentID,
	}
	tx := initializers.DB.Begin()
	defer tx.Rollback()

	nextTeamName := 0
	nextTeamNameS := ""
	result := tx.Where("tournament_id = ? and category_id = ?", body.TournamentID, body.CategoryID).Order("name DESC").First(&tournamentTeam)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		nextTeamNameS = "Team - 00"

	} else if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo buscando ultima pareja... ",
		})
		return
	} else {
		nextTeamNameS = strings.Split(tournamentTeam.Name, " ")[2]
		nextTeamName, _ = strconv.Atoi(nextTeamNameS)

	}
	nextTeamName++

	p1 := models.PlayersByTournament{
		Name:         fmt.Sprintf("Pareja - %02d, Jugador - 1", nextTeamName),
		Ranking:      body.Ranking1,
		CategoryID:   body.CategoryID,
		TournamentID: body.TournamentID,
	}
	result = tx.Create(&p1)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al crear jugador 1... ",
		})
		return
	}

	p2 := models.PlayersByTournament{
		Name:         fmt.Sprintf("Pareja - %02d, Jugador - 2", nextTeamName),
		Ranking:      body.Ranking1,
		CategoryID:   body.CategoryID,
		TournamentID: body.TournamentID,
	}
	result = tx.Create(&p2)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al crear jugador 2... ",
		})
		return
	}

	newTeam := models.TournamentTeam{
		Name:         fmt.Sprintf("Pareja - %02d", nextTeamName),
		CategoryID:   body.CategoryID,
		TournamentID: body.TournamentID,
		Player1ID:    p1.ID,
		Player2ID:    p2.ID,
	}
	result = tx.Create(&newTeam)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al crear pareja... ",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"data": newTeam})
}

func PutEnrolledTeams(c *gin.Context) {
	var body models.TournamentTeam

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al leer body...",
		})
		return
	}

	tournamentTeam := models.TournamentTeam{
		Name:         body.Name,
		PlayLevel:    body.PlayLevel,
		Player1ID:    body.Player1ID,
		Name1:        body.Name1,
		Ranking1:     body.Ranking1,
		Player2ID:    body.Player2ID,
		Name2:        body.Name2,
		Ranking2:     body.Ranking2,
		CategoryID:   body.CategoryID,
		TournamentID: body.TournamentID,
	}
	tournamentTeam.ID = body.ID

	result := initializers.DB.Save(&tournamentTeam)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al crear usuario... ",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ID": tournamentTeam.ID})
}

func GetEnrolledTeams(c *gin.Context) {
	var TournamentID = c.DefaultQuery("TournamentID", "0")
	var CategoryID = c.DefaultQuery("CategoryID", "0")

	type TeamExtended struct {
		models.TournamentTeam
		Paid1 int
		Paid2 int
	}
	var teamExtendedResults []TeamExtended

	sqlStatement := `Select distinct pdt.*, pdt.id as team_id, pr1.id as paid1, pr2.id as paid2
					from tournament_teams pdt 
					left join payment_records pr1 on pdt.member1_id = pr1.player_id
					left join payment_records pr2 on pdt.member2_id = pr2.player_id
					where pdt.tournament_id = ` + TournamentID + ` AND pdt.category_id= ` + CategoryID
	results := initializers.DB.Debug().Raw(sqlStatement).Scan(&teamExtendedResults)
	if results.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "results": len(teamExtendedResults), "data": teamExtendedResults})

}

func DeleteEnrolledTeam(c *gin.Context) {
	id := c.Param("id")

	tx := initializers.DB.Begin()
	defer tx.Rollback()

	var team models.TournamentTeam
	result := tx.Where("id = ?", id).First(&team)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fallo al eliminar... "})
		return
	}
	if team.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El equipo no existe... "})
		return
	}

	var p1 models.PlayersByTournament
	p1.ID = team.Player1ID
	result = tx.Unscoped().Delete(&p1)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fallo al eliminar jugador 1... "})
		return
	}

	var p2 models.PlayersByTournament
	p2.ID = team.Player2ID
	result = tx.Unscoped().Delete(&p2)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fallo al eliminar jugador 2... "})
		return
	}

	result = tx.Unscoped().Delete(&team)

	result = tx.Commit()
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fallo al eliminar pareja... "})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ID #" + id: "deleted"})
}
