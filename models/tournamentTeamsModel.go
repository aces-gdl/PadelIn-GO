package models

import (
	"gorm.io/gorm"
)

type TournamentTeam struct {
	gorm.Model
	Name           string
	PlayLevel      uint
	CategoryID     uint
	TeamID         uint
	Player1ID      uint
	Name1          string
	Reference1     string
	Restriction1   string
	FirstLastName1 string
	Ranking1       int
	Player2ID      uint
	Name2          string
	Reference2     string
	Restriction2   string
	FirstLastName2 string
	Ranking2       int
	TournamentID   uint
}
