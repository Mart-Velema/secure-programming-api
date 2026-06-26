package items

import (
	"strings"

	"guineatrade.nhlstenden.com/src/backpack"
)

type InventoryResponse struct {
	Success      int           `json:"success"`
	Descriptions []Description `json:"descriptions"`
	Assets       []Asset       `json:"assets"`
}

func (r *InventoryResponse) ToItem() Items {
	descriptions := make(map[string]*Description)
	for _, description := range r.Descriptions {
		descriptions[description.ClassId] = &description
	}
	itemResult := Items{}
	for _, asset := range r.Assets {
		description := descriptions[asset.ClassId]

		if description.Tradable != 1 {
			continue
		}

		defindex := backpack.GetDefindex(description.MarketHashName)
		if defindex == 0 {
			continue
		}
		marketHashName := backpack.GetMarketHashName(defindex)

		item := Item{
			AssetId:    asset.AssetId,
			InstanceId: asset.InstanceId,
			ClassId:    asset.ClassId,
			Defindex:   defindex,
			ItemType: ItemType{
				MarketHashName: marketHashName,
				Craftable:      description.getCraftability(),
				Quality:        description.getType(),
			},
		}

		effect, hasEffect := description.getUnusual()
		if hasEffect {
			item.Effect = effect
		}

		itemResult.Assets = append(itemResult.Assets, item)
	}

	return itemResult
}

type Description struct {
	Name           string             `json:"name"`
	ClassId        string             `json:"classid"`
	InstanceId     string             `json:"instanceid"`
	MarketHashName string             `json:"market_hash_name"`
	Tradable       int8               `json:"tradable"`
	Tags           []Tag              `json:"tags"`
	Description    []MetaDescriptions `json:"descriptions"`
}

type MetaDescriptions struct {
	Value string `json:"value"`
}

type Tag struct {
	Category     string `json:"category"`
	InternalName string `json:"internal_name"`
}

func (d *Description) getType() backpack.Quality {
	for _, tag := range d.Tags {
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

func (d *Description) getCraftability() bool {
	for _, description := range d.Description {
		if strings.Contains(description.Value, "Craft") {
			return false
		}
	}

	return true
}

func (d *Description) getUnusual() (string, bool) {
	for _, description := range d.Description {
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

type Asset struct {
	AssetId    string `json:"assetid"`
	ClassId    string `json:"classid"`
	InstanceId string `json:"instanceid"`
}
