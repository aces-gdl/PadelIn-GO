package initializers

import "PadelIn/models"

func SyncDatabase() {
	DB.AutoMigrate(
		&models.User{},
		&models.Institution{},
		&models.Origin{},
		&models.Payment{},
		&models.Tournament{},
		&models.TournamentCategory{},
		&models.Club{},
	)
}
