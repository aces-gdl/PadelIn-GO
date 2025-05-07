package models

import (
	"gorm.io/gorm"
)

type PlayersByTournament struct {
	gorm.Model    `gorm:"embedded"`
	TournamentID  uint
	CategoryID    uint
	Name          string
	Reference     string
	Ranking       int
	UserID        uint
	Restriction   string
	PaymentStatus string
}
