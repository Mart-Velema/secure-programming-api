package items

import "guineatrade.nhlstenden.com/src/backpack"

type Items struct {
	Assets []Item
}

type Item struct {
	AssetId    string `json:"assetid"`
	InstanceId string `json:"instanceid"`
	ClassId    string `json:"classid"`
	Defindex   uint32 `json:"defindex"`
	ItemType
}

type ItemType struct {
	MarketHashName string           `json:"marketHashName"`
	Craftable      bool             `json:"craftability"`
	Quality        backpack.Quality `json:"quality"`
	Effect         string           `json:"unusual,omitempty"`
}

type Stock struct {
	ItemType
	Quantity uint32 `json:"quantity"`
}

func (items Items) ToStock() []Stock {
	stockMap := make(map[ItemType]uint32)

	for _, item := range items.Assets {
		stockMap[item.ItemType]++
	}

	var stock []Stock
	for itemType, quantity := range stockMap {
		stock = append(stock, Stock{itemType, quantity})
	}

	return stock
}

func (items Items) GetItemsOfType(itemType ItemType, quantity uint32) []Item {
	itemList := make([]Item, quantity)
	for _, item := range items.Assets {
		if item.ItemType == itemType {
			itemList = append(itemList, item)
			if uint32(len(itemList)) == quantity {
				return itemList
			}
		}
	}
	return nil
}
