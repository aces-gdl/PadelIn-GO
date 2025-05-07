package controllers

import (
	initializers "PadelIn/initializers"
	models "PadelIn/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func PostRegisterUserForTournament(c *gin.Context) {
	var body models.PlayersByTournament

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al leer body...",
		})
		return
	}

	if body.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Debe especificar un usuario final...",
		})
		return
	}

	if body.TournamentID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Debe especificar un torneo...",
		})
		return
	}
	if body.CategoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Debe especificar una categoria...",
		})
		return
	}

	var tournament models.Tournament
	tournament.ID = body.TournamentID
	result := initializers.DB.Debug().Where(&tournament).First(&tournament)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Torneo no encontrado...",
		})
		return
	}

	var category models.TournamentCategory
	category.ID = body.CategoryID
	result = initializers.DB.Debug().Where(&category).First(&category)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Categoría no encontrada...",
		})
		return
	}

	var endUser models.EndUser
	endUser.ID = body.UserID
	result = initializers.DB.Debug().Where(&endUser).First(&endUser)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usuario final no encontrado...",
		})
		return
	}

	playerForTournament := models.PlayersByTournament{
		TournamentID:  body.TournamentID,
		CategoryID:    body.CategoryID,
		UserID:        body.UserID,
		PaymentStatus: "Pendiente",
		Name:          endUser.Name,
	}

	initializers.DB.Debug().Create(&playerForTournament)

	c.JSON(http.StatusOK, gin.H{"ID": playerForTournament.ID})
}

func DeleteRegisterUserForTournament(c *gin.Context) {
	TournamentID := c.DefaultQuery("TournamentID", "")
	CategoryID := c.DefaultQuery("CategoryID", "")
	UserID := c.DefaultQuery("UserID", "")

	if TournamentID == "" || CategoryID == "" || UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al leer parametros...",
		})
		return
	}

	var playerForTournament models.PlayersByTournament

	statement := `delete from players_by_tournaments where  tournament_id = ` + TournamentID +
		` and category_id = ` + CategoryID +
		` and user_id = ` + UserID
	results := initializers.DB.Debug().Exec(statement)

	c.JSON(http.StatusOK, gin.H{"ID": playerForTournament.ID, "count": results.RowsAffected})
}

func GetRegisteredUsersByTournament(c *gin.Context) {
	var categoryID = c.DefaultQuery("CategoryID", "")
	var tournamentID = c.DefaultQuery("TournamentID", "")
	var SearchString = c.DefaultQuery("SearchString", "")
	var status = c.DefaultQuery("Status", "")

	var whereClause = " where 1 = 1 "
	var whereClauseSearchString = ""
	var whereClauseStatus = ""
	if categoryID != "" {
		whereClause = whereClause + " AND u.category_id = " + categoryID
	}
	if status != "" {
		if status == "2" {
			whereClauseStatus = whereClauseStatus + " AND pbt.payment_status IS NOT NULL "
		}
		if status == "3" {
			whereClauseStatus = whereClauseStatus + " AND pbt.payment_status IS NULL "
		}
	}
	if SearchString != "" {
		whereClauseSearchString = " AND u.name like  '%" + SearchString + "%'"
	}

	queryString := `SELECT u.*, 
		c.description as category_description, 
		c.color as category_color,
		pbt.payment_status 
	FROM "users" u
		inner join categories c on u.category_id = c.id 
		left join players_by_tournaments pbt on u.id = pbt.user_id  and c.id = pbt.category_id and pbt.tournament_id = ` + tournamentID + `
 ` + whereClause + whereClauseSearchString + whereClauseStatus + `
		Order by  u.name asc`

	type userExtended struct {
		ID                    uint
		CategoryID            uint
		CategoryColor         string
		Name                  string
		FamilyName            string
		GivenName             string
		Email                 string
		CategoryDescription   string
		PermissionID          uint
		Phone                 string
		HasPicture            int
		MemberSince           time.Time
		Birthday              time.Time
		PermissionDescription string
		Ranking               int
		PaymentStatus         string
	}

	var usersExtended []userExtended

	results := initializers.DB.Debug().Raw(queryString).Where(whereClause).Scan(&usersExtended)
	if results.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	type stats struct {
		Inscritos   uint
		NoInscritos uint
		Todos       uint
	}

	var statsRecord stats
	getStatsStament := `SELECT SUM(case when pbt.payment_status IS NOT NULL then 1 else 0 end ) as inscritos ,
							   SUM(case when pbt.payment_status IS NULL then 1 else 0 end ) as no_inscritos ,
							   count(* )as todos
						FROM "users" u
								inner join categories c on u.category_id = c.id 
								left join players_by_tournaments pbt on u.id = pbt.user_id  and c.id = pbt.category_id 	and pbt.tournament_id =	 ` + tournamentID + `
						` + whereClause
	initializers.DB.Debug().Raw(getStatsStament).Scan(&statsRecord)
	c.JSON(http.StatusOK, gin.H{"status": "success", "results": len(usersExtended), "data": usersExtended, "stats": statsRecord})
}

func UpdatePlayerByTournament(c *gin.Context) {
	// Get the player ID from the URL parameter
	playerID := c.Param("id")

	var body models.PlayersByTournament

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the existing player by tournament record
	var playerByTournament models.PlayersByTournament
	if err := initializers.DB.First(&playerByTournament, playerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player by tournament record not found"})
		return
	}

	playerByTournament.Name = body.Name
	playerByTournament.Reference = body.Reference
	playerByTournament.Ranking = body.Ranking
	playerByTournament.Restriction = body.Restriction

	// Save the updated record
	if err := initializers.DB.Save(&playerByTournament).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update player by tournament"})
		return
	}

	// Return the updated record
	c.JSON(http.StatusOK, gin.H{"data": playerByTournament})
}
