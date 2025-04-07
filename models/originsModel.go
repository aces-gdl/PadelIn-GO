package models

import (
	"gorm.io/gorm"
)

type Origin struct {
	gorm.Model
	URL           string
	InstitutionID uint
	Active        bool
}
