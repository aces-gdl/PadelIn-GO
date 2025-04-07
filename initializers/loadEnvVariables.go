package initializers

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	godotenv.Load()
	/*
		if err != nil {
			log.("Error cargando archivo .env")
		} */
	fmt.Println("...Inicializando")

}
