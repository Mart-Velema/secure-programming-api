package _nit

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	envFileLocation := "./.env"
	for idx, arg := range os.Args {
		if arg == "--env" {
			envFileLocation = os.Args[idx+1]
		}
	}
	log.Printf("Loading .env file from: %s", envFileLocation)
	err := godotenv.Load(envFileLocation)
	if err != nil {
		log.Fatalf("Error loading .env file: %s\n", err)
	}
}
