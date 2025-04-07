package models

import (
	"time"

	"gorm.io/gorm"
)

type Tournament struct {
	gorm.Model
	Description      string
	ClubName         string
	StartDate        time.Time
	EndDate          time.Time
	StartTime        time.Time
	EndTime          time.Time
	HostClubID       uint
	GameDuration     int
	RoundrobinCourts int
	IdInstitution    int
	Active           bool
}
