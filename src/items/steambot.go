package items

import (
	"strings"

	"guineatrade.nhlstenden.com/src/backpack"
)

type SteamInventoryResponse struct {
	Inventory []SteamItem `json:"inventory"`
}

func (s *SteamInventoryResponse) ToItems() Items {
	itemResult := Items{
		Assets: make([]Item, len(s.Inventory)),
	}
	for idx, asset := range s.Inventory {
		itemResult.Assets[idx] = Item{
			AssetId:        asset.AssetId,
			InstanceId:     asset.InstanceId,
			ClassId:        asset.ClassId,
			MarketHashName: asset.MarketHashName,
			Craftable:      asset.getCraftability(),
			Quality:        asset.getType(),
			Effect:         "Whirly Wind",
		}
	}

	return itemResult
}

type SteamItem struct {
	AssetId        string             `json:"assetid"`
	InstanceId     string             `json:"instanceid"`
	ClassId        string             `json:"classid"`
	MarketHashName string             `json:"market_hash_name"`
	Tradable       bool               `json:"tradable"`
	Descriptions   []MetaDescriptions `json:"descriptions"`
	Tags           []Tag              `json:"tags"`
}

func (s *SteamItem) getType() backpack.Quality {
	for _, tag := range s.Tags {
		if tag.Category != "Quality" {
			continue
		}

		switch tag.InternalName {
		case "Rarity4":
			return backpack.Unusual
		case "Unique":
			return backpack.Unique
		case "strange":
			return backpack.Strange

		}
	}

	return backpack.Unique
}

func (s *SteamItem) getCraftability() bool {
	for _, description := range s.Descriptions {
		if strings.Contains(description.Value, "Craft") {
			return false
		}
	}

	return true
}
