package models

import "gorm.io/gorm"

type Institution struct {
	gorm.Model
	Name        string
	Description string
	Origin      string
}
