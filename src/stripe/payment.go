package stripe

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v85"
	"gorm.io/gorm"
	"guineatrade.nhlstenden.com/src/auth/middleware"
	"guineatrade.nhlstenden.com/src/database"
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
	var discount int64
	for i := range checkoutItems {
		var checkoutItem = checkoutItems[i]
		var price = int64(checkoutItem.Price())

		if checkoutItem.IsSold {
			totalCost -= price
			discount += price
		} else {
			totalCost += price
		}
	}
	discount *= 2 // Output of function is all prices of sold items combined *2 to get the discount instead

	if totalCost > 0 {
		assets := toAssets(checkoutItems)

		trade := database.Trade{
			UserID:      user.ID,
			Cost:        totalCost,
			TradeAction: database.BUY,
			TradeStatus: database.PAYMENT_IN_PROGRESS,
			Assets:      assets,
		}

		err = database.GetInstance().Session(&gorm.Session{FullSaveAssociations: true}).Create(&trade).Error

		if err != nil {
			log.Println("ERROR WITH CREATING TRADE IN DB", err)
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Could not create transaction"})
			return
		}

		lineItems := make([]*stripe.CheckoutSessionCreateLineItemParams, 0, len(requestedStockList))

		for _, item := range checkoutItems {
			lineItems = append(lineItems, createCheckoutSessionLineItem(
				int64(item.Price()),
				item.Name(),
				item.Description(),
				int64(len(item.Items)),
			))
		}

		session, err := createPaymentSessionData(lineItems, discount, strconv.Itoa(int(trade.ID)))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create payment session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "created_link", "url": session.URL})
	} else {
		res, err := sendCheckoutTradeOffer(user, checkoutItems)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create trade offer"})
		}

		trade := database.Trade{
			UserID:       user.ID,
			Cost:         totalCost,
			TradeAction:  database.SELL,
			TradeStatus:  database.TRADE_IN_PROGRESS,
			SteamTradeId: res.TradeOfferID,
		}

		database.GetInstance().Save(&trade)

		c.JSON(http.StatusOK, gin.H{"status": "no_payment_required", "receives": totalCost})
	}
}

func sendTradeOffer(user database.User, sellItems []steam.TradeOfferItem, buyItems []steam.TradeOfferItem) (*steam.SendTradeOfferResponse, error) {
	req := steam.SendTradeOfferRequest{
		TradeURL:       user.TradeUrl,
		ItemsToGive:    buyItems,
		ItemsToReceive: sellItems,
		Message:        "Thanks for trading with GuineaTrade!",
	}

	res, err := steam.SendTradeOffer(req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func sendCheckoutTradeOffer(user database.User, checkoutItems []CheckoutItem) (*steam.SendTradeOfferResponse, error) {
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

	return sendTradeOffer(user, sellItems, buyItems)
}

func createCheckoutSessionLineItem(price int64, name string, description string, quantity int64) *stripe.CheckoutSessionCreateLineItemParams {
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
	}
}

func createPaymentSessionData(lineItems []*stripe.CheckoutSessionCreateLineItemParams, discount int64, transactionId string) (*stripe.CheckoutSession, error) {

	params := &stripe.CheckoutSessionCreateParams{
		SuccessURL: stripe.String("https://google.com/success"),
		CancelURL:  stripe.String("https://google.com/cancel"),
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems:  lineItems,
		Metadata: map[string]string{
			"transaction_id": transactionId,
		},
	}

	if discount > 0 {
		couponParams := &stripe.CouponCreateParams{
			Name:      new("Sold Items"),
			AmountOff: new(discount),
			Currency:  stripe.String(string(stripe.CurrencyUSD)),
			Duration:  stripe.String(string(stripe.CouponDurationOnce)),
		}

		c, err := sc.V1Coupons.Create(context.TODO(), couponParams)
		if err != nil {
			return nil, err
		}

		params.Discounts = []*stripe.CheckoutSessionCreateDiscountParams{
			{
				Coupon: stripe.String(c.ID),
			},
		}
	}

	return sc.V1CheckoutSessions.Create(context.TODO(), params)
}
