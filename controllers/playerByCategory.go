package controllers

import (
	"PadelIn/initializers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PlayersByCategory(c *gin.Context) {

	type PlayerByCategory struct {
		ID                  uint
		Name                string
		Phone               string
		CategoryID          int
		CategoryDescription string
	}

	sqlStatement := "select "
	sqlStatement += " eu.id, eu.name,eu.phone, "
	sqlStatement += " case when pbt.category_id IS NULL then 0 else pbt.category_id end as category_id, "
	sqlStatement += " case when pbt.category_id IS NULL then 'N/D' else tc.description end as category_description "
	sqlStatement += " from end_users eu "
	sqlStatement += " left join players_by_tournaments pbt on eu.id = pbt.user_id "
	sqlStatement += " left join tournament_categories tc on pbt.category_id = tc.id"

	var PlayersByCategory []PlayerByCategory
	result := initializers.DB.Raw(sqlStatement).Scan(&PlayersByCategory)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "No existen jugadores..."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "results": len(PlayersByCategory), "data": PlayersByCategory})
}
