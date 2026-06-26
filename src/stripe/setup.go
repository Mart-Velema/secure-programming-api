package stripe

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v85"
	"guineatrade.nhlstenden.com/src/backpack"
	"guineatrade.nhlstenden.com/src/database"
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

func (c CheckoutItem) Name() string {
	return c.Items[0].MarketHashName
}

func (c CheckoutItem) Description() string {
	item := c.Items[0]
	craftable := ""
	if item.Craftable {
		craftable = "Craftable"
	} else {
		craftable = "Uncraftable"
	}

	return fmt.Sprintf("%s - %s", craftable, item.Quality)
}

func (c CheckoutItem) Price() uint {
	return uint(len(c.Items)) * c.PricePerItem
}

type CheckoutRequest struct {
	MarketHashName string           `json:"marketHashName"`
	Craftable      bool             `json:"craftability"`
	Quality        backpack.Quality `json:"quality"`
	Effect         string           `json:"unusual,omitempty"`
	Quantity       uint32           `json:"quantity"`
	IsSold         bool             `json:"isSold"`
}

type Checkout struct {
	stock  items.Stock
	isSold bool
}

func (cr CheckoutRequest) ToCheckout() Checkout {
	return Checkout{
		stock: items.Stock{
			ItemType: items.ItemType{
				MarketHashName: cr.MarketHashName,
				Craftable:      cr.Craftable,
				Quality:        cr.Quality,
				Effect:         cr.Effect,
			},
			Quantity: cr.Quantity,
		},
		isSold: cr.IsSold,
	}
}

func toAssets(checkoutItems []CheckoutItem) []database.Asset {
	var assets []database.Asset
	for _, checkoutItem := range checkoutItems {
		for _, item := range checkoutItem.Items {
			var asset = database.Asset{
				AssetId: item.AssetId,
			}

			if checkoutItem.IsSold {
				asset.TradeDirection = database.SELL
			} else {
				asset.TradeDirection = database.BUY
			}
			assets = append(assets, asset)
		}
	}
	return assets
}

func sortAssetsToTradeOfferItems(assets []database.Asset) ([]steam.TradeOfferItem, []steam.TradeOfferItem) {
	var sellItems []database.Asset
	var buyItems []database.Asset

	for _, item := range assets {
		if item.TradeDirection == database.SELL {
			sellItems = append(sellItems, item)
		} else {
			buyItems = append(buyItems, item)
		}
	}

	return toTradeOfferItems(sellItems), toTradeOfferItems(buyItems)
}

func toTradeOfferItems(itemAssets []database.Asset) []steam.TradeOfferItem {
	var tradeOfferItems []steam.TradeOfferItem
	for _, asset := range itemAssets {
		tradeOfferItems = append(tradeOfferItems, steam.TradeOfferItem{
			AppID:     440,
			ContextID: "2",
			AssetID:   asset.AssetId,
		})
	}

	return tradeOfferItems
}

func toCheckoutItems(steamId uint64, stockList []CheckoutRequest) ([]CheckoutItem, error) {
	groups := make(map[items.ItemType]*CheckoutItem)

	userInv, err := inventory.GetUserInventory(steamId)
	if err != nil {
		return nil, err
	}
	var userInventory = userInv.ToItem()

	botInventory, err := steam.GetBotInventoryData()
	if err != nil {
		return nil, err
	}

	for _, item := range stockList {
		checkoutItem := item.ToCheckout()
		stock := checkoutItem.stock

		var itemInventory items.Items
		if item.IsSold {
			itemInventory = userInventory
		} else {
			itemInventory = botInventory
		}

		assetItems := itemInventory.GetItemsOfType(stock.ItemType, stock.Quantity)
		if assetItems == nil {
			return nil, errors.New("could not get stock")
		}

		groups[stock.ItemType] = &CheckoutItem{
			Items:        assetItems,
			PricePerItem: backpack.GetSpecificPrice(stock.MarketHashName, stock.Quality, stock.Craftable, stock.Effect),
			IsSold:       item.IsSold,
		}
	}

	result := make([]CheckoutItem, 0, len(groups))
	for _, checkoutItem := range groups {
		result = append(result, *checkoutItem)
	}

	return result, nil
}
