package stripe

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v85"
)

var (
	sc *stripe.Client
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s\n", err)
	}

	stripeApiKey, apiKeyExists := os.LookupEnv("STRIPE_SECRET_KEY")

	if !apiKeyExists {
		log.Fatal("STRIPE_SECRET_KEY is unset")
	}

	sc = stripe.NewClient(stripeApiKey)
}
