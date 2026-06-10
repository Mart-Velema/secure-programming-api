package inventory

type InventoryResponse struct {
	Success      int           `json:"success"`
	Descriptions []Description `json:"descriptions"`
	Assets       []Asset       `json:"assets"`
}

type Description struct {
	Name       string `json:"name"`
	IconUrl    string `json:"icon_url"`
	ClassId    string `json:"classid"`
	InstanceId string `json:"instanceid"`
}

type Asset struct {
	AssetId    string `json:"assetid"`
	ClassId    string `json:"classid"`
	InstanceId string `json:"instanceid"`
}
