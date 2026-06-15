package items

import (
	"fmt"
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
		if asset.MarketHashName == "Unusual Taunt: Square Dance" {
			fmt.Println(asset.MarketHashName)
		}
		defindex := backpack.GetDefindex(asset.MarketHashName)
		if defindex == 0 {
			continue
		}
		marketHashName := backpack.GetMarketHashName(defindex)

		item := Item{
			AssetId:        asset.AssetId,
			InstanceId:     asset.InstanceId,
			ClassId:        asset.ClassId,
			MarketHashName: marketHashName,
			Defindex:       defindex,
			Craftable:      asset.getCraftability(),
			Quality:        asset.getType(),
		}

		effect, hasEffect := asset.getUnusual()
		if hasEffect {
			item.Effect = effect
		}

		itemResult.Assets[idx] = item
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
		case "rarity4":
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

func (s SteamItem) getUnusual() (string, bool) {
	for _, description := range s.Descriptions {
		if strings.Contains(description.Value, "Unusual") {
			splits := strings.Split(description.Value, ":")
			if len(splits) != 2 {
				continue
			}

			return strings.TrimSpace(splits[1]), true
		}
	}

	return "", false
}
