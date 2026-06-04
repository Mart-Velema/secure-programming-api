package stripe

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v85"
)

func CreatePaymentSession(c *gin.Context) {
	lineItems := []*stripe.CheckoutSessionCreateLineItemParams{
		{
			PriceData: &stripe.CheckoutSessionCreateLineItemPriceDataParams{
				Currency:   stripe.String("usd"),
				UnitAmount: stripe.Int64(1000),
				ProductData: &stripe.CheckoutSessionCreateLineItemPriceDataProductDataParams{
					Name:        stripe.String("Item A"),
					Description: stripe.String("First item"),
				},
			},
			Quantity: stripe.Int64(1),
		},
		{
			PriceData: &stripe.CheckoutSessionCreateLineItemPriceDataParams{
				Currency:   stripe.String("usd"),
				UnitAmount: stripe.Int64(2500),
				ProductData: &stripe.CheckoutSessionCreateLineItemPriceDataProductDataParams{
					Name:        stripe.String("Item B"),
					Description: stripe.String("Second item"),
				},
			},
			Quantity: stripe.Int64(2),
		},
		{
			PriceData: &stripe.CheckoutSessionCreateLineItemPriceDataParams{
				Currency:   stripe.String("usd"),
				UnitAmount: stripe.Int64(2500), // $25.00
				ProductData: &stripe.CheckoutSessionCreateLineItemPriceDataProductDataParams{
					Name:        stripe.String("Item C"),
					Description: stripe.String("Second item"),
				},
			},
			Quantity: stripe.Int64(1),
		},
	}
	session, err := createPaymentSessionData(lineItems)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	c.String(200, session.URL)
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
