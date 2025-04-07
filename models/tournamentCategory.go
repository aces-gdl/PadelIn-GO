package models

import (
	"gorm.io/gorm"
)

type TournamentCategory struct {
	gorm.Model

	Description   string
	BgColor       string
	TextColor     string
	Level         int
	Active        bool
	IdInstitution int
}
