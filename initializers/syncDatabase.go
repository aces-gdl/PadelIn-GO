package initializers

import "PadelIn/models"

func SyncDatabase() {
	DB.AutoMigrate(
		&models.EndUser{},
		&models.Institution{},
		&models.Origin{},
		&models.Payment{},
		&models.Tournament{},
		&models.TournamentCategory{},
		&models.Club{},
		&models.PlayersByTournament{},
		&models.TournamentTeam{},
	)
}
