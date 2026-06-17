package stripe

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v85"
	"guineatrade.nhlstenden.com/src/backpack"
	"guineatrade.nhlstenden.com/src/inventory"
	"guineatrade.nhlstenden.com/src/items"
	"guineatrade.nhlstenden.com/src/steam"
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

type CheckoutItem struct {
	Items        []items.Item
	PricePerItem uint
	IsSold       bool // Whether the item is bought or sold by the user
}

type CheckoutRequest struct {
	Stock  items.Stock
	IsSold bool
}

func toCheckoutItems(steamId uint64, stockList []CheckoutRequest) []CheckoutItem {
	groups := make(map[items.ItemType]*CheckoutItem)

	for _, item := range stockList {
		stock := item.Stock

		var itemInventory items.Items
		if item.IsSold {
			var err error
			if itemInventory, err = steam.GetBotInventoryData(); err != nil {
				return nil
			}
		} else {
			inv, err := inventory.GetUserInventory(steamId)
			if err != nil {
				return nil
			}
			itemInventory = inv.ToItem()
		}

		groups[stock.ItemType] = &CheckoutItem{
			Items:        itemInventory.GetItemsOfType(stock.ItemType, stock.Quantity), // TODO: Call this function safely
			PricePerItem: backpack.GetSpecificPrice(stock.MarketHashName, stock.Quality, stock.Craftable, stock.Effect),
			IsSold:       item.IsSold,
		}
	}

	result := make([]CheckoutItem, 0, len(groups))
	for _, checkoutItem := range groups {
		result = append(result, *checkoutItem)
	}

	return result
}

func (c CheckoutItem) Name() string {
	return c.Items[0].MarketHashName
}

func (c CheckoutItem) Description() string {
	item := c.Items[0]
	craftable := ""
	if item.Craftable == true {
		craftable = "Craftable"
	} else {
		craftable = "Uncraftable"
	}

	return fmt.Sprintf("%s - %s", craftable, item.Quality)
}

func (c CheckoutItem) Price() uint {
	return uint(len(c.Items)) * c.PricePerItem
}
