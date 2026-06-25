package stripe

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v85"
	"guineatrade.nhlstenden.com/src/database"
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
		log.Println("Error reading payload:", err)
		c.String(400, "Invalid Payload")
		return
	}

	evt, err := sc.ConstructEvent(payload, c.GetHeader("Stripe-Signature"), whSec)
	if err != nil {
		log.Println("Error verifying webhook signature:", err)
		c.String(400, "Webhook signature verification failed")
		return
	}

	switch evt.Type {
	case "checkout.session.completed":
		var checkoutSession stripe.CheckoutSession
		err := json.Unmarshal(evt.Data.Raw, &checkoutSession)
		if err != nil {
			log.Printf("Error unmarshalling checkout session: %s\n", err)
			return
		}

		var trade database.Trade

		if result := database.GetInstance().
			Preload("Assets").
			Where("id = ?", checkoutSession.Metadata["transaction_id"]).
			First(&trade); result.Error != nil {
			log.Println("Error getting trade:", err)
			c.String(400, "Error getting trade")
			return
		}

		var user database.User
		err = database.GetInstance().
			Where("id = ?", trade.UserID).
			First(&user).
			Error

		if err != nil {
			log.Printf("Error getting trade user: %s\n", err)
			c.String(400, "Error getting trade")
			return
		}

		sellItems, buyItems := sortAssetsToTradeOfferItems(trade.Assets)

		res, err := sendTradeOffer(user, sellItems, buyItems)

		if err != nil {
			log.Println("Error sending trade offer:", err)
			c.String(400, "Error getting trade")
			return
		}

		trade.TradeStatus = database.TRADE_IN_PROGRESS
		trade.SteamTradeId = res.TradeOfferID

		database.GetInstance().Save(&trade)

		c.String(200, checkoutSession.Metadata["transaction-id"])
	}
}
