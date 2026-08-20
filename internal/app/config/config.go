package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/KDarenskii/catalog-service/internal/app/config/section"
)

type Config struct {
	Repository section.Repository
	Processor  section.Processor
	Monitor    section.Monitor
}

var Root Config

func Load() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Не удалось загрузить файл")
	}

	err = envconfig.Process("APP", &Root)
	if err != nil {
		log.Fatalf("config: %v", err)
	}
}
