package stripe

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v85"
)

var (
	whSec string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s\n", err)
	}

	webhookSecretKey, apiKeyExists := os.LookupEnv("STRIPE_WH_SECRET_KEY")

	if !apiKeyExists {
		log.Fatal("STRIPE_WH_SECRET_KEY is unset")
	}

	whSec = webhookSecretKey
}

func Webhook(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid Payload")
		log.Printf("Error reading payload: %s\n", err)
		return
	}

	evt, err := sc.ConstructEvent(payload, c.GetHeader("Stripe-Signature"), whSec)
	if err != nil {
		c.String(400, "Webhook signature verification failed")
		log.Printf("Error verifying webhook signature: %s\n\n", err)
		return
	}

	switch evt.Type {
	case "checkout.session.completed":
		var checkoutComplete stripe.CheckoutSession
		err := json.Unmarshal(evt.Data.Raw, &checkoutComplete)
		if err != nil {
			log.Printf("Error unmarshalling checkout session: %s\n", err)
			return
		}

		c.String(200, checkoutComplete.Metadata["order-id"])
		log.Println(checkoutComplete.Metadata["order-id"])
	}

	//body, err := io.ReadAll(c.Request.Body)
	//if err != nil {
	//	fmt.Println("error reading body:", err)
	//	return
	//}
	//
	//fmt.Println("Request body:")
	//fmt.Println(string(body))
}
