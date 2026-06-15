package inventory

import (
	"strings"

	"guineatrade.nhlstenden.com/src/backpack"
	"guineatrade.nhlstenden.com/src/items"
)

type InventoryResponse struct {
	Success      int           `json:"success"`
	Descriptions []Description `json:"descriptions"`
	Assets       []Asset       `json:"assets"`
}

func (r *InventoryResponse) toItem() items.Items {
	descriptions := make(map[string]*Description)
	for _, description := range r.Descriptions {
		descriptions[description.ClassId] = &description
	}
	itemResult := items.Items{
		Assets: make([]items.Item, len(r.Assets)),
	}
	for idx, asset := range r.Assets {
		description := descriptions[asset.ClassId]
		itemResult.Assets[idx] = items.Item{
			AssetId:        asset.AssetId,
			InstanceId:     asset.InstanceId,
			ClassId:        asset.ClassId,
			MarketHashName: description.MarketHashName,
			Craftable:      description.getCraftability(),
			Quality:        description.getType(),
			Effect:         "Whirly Wind",
		}
	}

	return itemResult
}

type Description struct {
	Name           string             `json:"name"`
	ClassId        string             `json:"classid"`
	InstanceId     string             `json:"instanceid"`
	MarketHashName string             `json:"market_hash_name"`
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

func (d *Description) getCraftability() bool {
	for _, description := range d.Description {
		if strings.Contains(description.Value, "Craft") {
			return false
		}
	}

	return true
}

type Asset struct {
	AssetId    string `json:"assetid"`
	ClassId    string `json:"classid"`
	InstanceId string `json:"instanceid"`
}
