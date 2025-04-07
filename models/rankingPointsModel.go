package models

import (
	"time"

	"gorm.io/gorm"
)

type RankingPoint struct {
	gorm.Model
	UserID         uint      `gorm:"not null;index"`
	InstitutionID  uint      `gorm:"index"`
	TournamentID   uint      `gorm:"index"`
	Points         int       `gorm:"not null"`
	Category       string    `gorm:"type:varchar(50);not null"`
	AwardDate      time.Time `gorm:"not null;index"`
	ExpirationDate time.Time `gorm:"index"`
	Reason         string    `gorm:"type:varchar(255)"`
	IsActive       bool      `gorm:"default:true"`
	Season         int       `gorm:"not null;index"`
	MetaData       string    `gorm:"type:json"`
}

type RankingPointsHistory struct {
	gorm.Model
	RankingPointsID uint      `gorm:"not null;index"`
	OldPoints       int       `gorm:"not null"`
	NewPoints       int       `gorm:"not null"`
	ChangeDate      time.Time `gorm:"not null;index"`
	Reason          string    `gorm:"type:varchar(255)"`
}
