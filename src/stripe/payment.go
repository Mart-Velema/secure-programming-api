package stripe

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v85"
	"guineatrade.nhlstenden.com/src/auth/middleware"
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

	lineItems := make([]*stripe.CheckoutSessionCreateLineItemParams, 0, len(requestedStockList))

	for _, item := range toCheckoutItems(user.SteamId, requestedStockList) {
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

	c.String(http.StatusOK, session.URL)
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
