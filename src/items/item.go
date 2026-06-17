package items

import "guineatrade.nhlstenden.com/src/backpack"

type Items struct {
	Assets []Item
}

type Item struct {
	AssetId        string           `json:"assetid"`
	InstanceId     string           `json:"instanceid"`
	ClassId        string           `json:"classid"`
	Defindex       uint32           `json:"defindex"`
	MarketHashName string           `json:"marketHashName"`
	Craftable      bool             `json:"craftability"`
	Quality        backpack.Quality `json:"quality"`
	Effect         string           `json:"unusual,omitempty"`
}
