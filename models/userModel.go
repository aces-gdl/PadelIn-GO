package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name          string
	Type          string
	Email         string
	User          string
	Phone         string `gorm:"not null;index:idx_phone_institution,priority:1,unique"`
	Password      string
	DateOfBirth   time.Time
	InstitutionID uint `gorm:"not null;index:idx_phone_institution,priority:2,unique"`
	Active        string
	PassCode      string
}
