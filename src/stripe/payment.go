package stripe

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v85"
	"guineatrade.nhlstenden.com/src/auth/middleware"
	"guineatrade.nhlstenden.com/src/steam"
)

func CreatePaymentSession(c *gin.Context) {
	var requestedStockList []CheckoutRequest
	if err := c.ShouldBindJSON(&requestedStockList); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid json input"})
		return
	}

	user, err := middleware.ExtractTokenUser(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Token expired"})
		return
	}

	checkoutItems, err := toCheckoutItems(user.SteamId, requestedStockList)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Could not get items"})
		return
	}

	var totalCost int64
	for i := range checkoutItems {
		var checkoutItem = checkoutItems[i]
		var price = int64(checkoutItem.Price())

		if checkoutItem.IsSold {
			totalCost -= price
		} else {
			totalCost += price
		}
	}

	if totalCost > 0 {
		lineItems := make([]*stripe.CheckoutSessionCreateLineItemParams, 0, len(requestedStockList))

		for _, item := range checkoutItems {
			lineItems = append(lineItems, createCheckoutSessionLineItem(
				int64(item.Price()),
				item.Name(),
				item.Description(),
				int64(len(item.Items)),
				"", // TODO: Add transactionIds
			))
		}

		session, err := createPaymentSessionData(lineItems)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create payment session"})
			return
		}

		// TODO: Save a record of this in the database
		c.JSON(http.StatusOK, gin.H{"status": "created_link", "url": session.URL})
	} else {
		var sellItems []steam.TradeOfferItem
		var buyItems []steam.TradeOfferItem
		for _, checkoutItem := range checkoutItems {
			for _, item := range checkoutItem.Items {
				tradeOfferItem := steam.TradeOfferItem{
					AppID:     440,
					ContextID: "2",
					AssetID:   item.AssetId,
				}

				if checkoutItem.IsSold {
					sellItems = append(sellItems, tradeOfferItem)
				} else {
					buyItems = append(buyItems, tradeOfferItem)
				}
			}
		}

		req := steam.SendTradeOfferRequest{
			TradeURL:       user.TradeUrl,
			ItemsToGive:    buyItems,
			ItemsToReceive: sellItems,
			Message:        "Thanks for trading with GuineaTrade!",
		}

		err = steam.SendTradeOffer(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create trade offer"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "no_payment_required", "receives": totalCost})
	}
}

func createCheckoutSessionLineItem(price int64, name string, description string, quantity int64, transactionId string) *stripe.CheckoutSessionCreateLineItemParams {
	return &stripe.CheckoutSessionCreateLineItemParams{
		PriceData: &stripe.CheckoutSessionCreateLineItemPriceDataParams{
			Currency:   new("usd"),
			UnitAmount: new(price),
			ProductData: &stripe.CheckoutSessionCreateLineItemPriceDataProductDataParams{
				Name:        new(name),
				Description: new(description),
			},
		},
		Quantity: new(quantity),
		Metadata: map[string]string{
			"transactionId": transactionId,
		},
	}
}

func createPaymentSessionData(lineItems []*stripe.CheckoutSessionCreateLineItemParams) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionCreateParams{
		SuccessURL: stripe.String("https://google.com/success"),
		CancelURL:  stripe.String("https://google.com/cancel"),
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems:  lineItems,
		Metadata: map[string]string{
			"order_id": "1234",
		},
	}

	return sc.V1CheckoutSessions.Create(context.TODO(), params)
}
