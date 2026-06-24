package _nit

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	envFileLocation, found := os.LookupEnv("ENV_FILE")
	if !found {
		envFileLocation = "./.env"
		log.Println("Using default path for .env file: ./.env")
	}
	log.Printf("Loading .env file from: %s", envFileLocation)
	err := godotenv.Load(envFileLocation)
	if err != nil {
		log.Fatalf("Error loading .env file: %s\n", err)
	}
}
