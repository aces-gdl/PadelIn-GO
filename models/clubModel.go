package models

import (
	"gorm.io/gorm"
)

type Club struct {
	gorm.Model
	Name          string
	Description   string
	Contact       string
	ImageURL      string
	Address       string
	Phone         string
	InstitutionID int
}
