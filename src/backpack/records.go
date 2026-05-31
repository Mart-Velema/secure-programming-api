package backpack

import (
	"log"
	"strconv"
	"time"
)

type PricingItems string

const (
	PricingMetal = "metal"
	PricingKey   = "keys"
)

type QualityItems int

const (
	Normal QualityItems = iota
	Genuine
	Rarity2
	Vintage
	Rarity3
	Unusual
	Unique
	Community
	Valve
	SelfMade
	Customized
	Strange
	Completed
	Haunted
	Collectors
	Decorated
)

var qualityMap = map[string]QualityItems{
	"0":  Normal,
	"1":  Genuine,
	"2":  Rarity2,
	"3":  Vintage,
	"4":  Rarity3,
	"5":  Unusual,
	"6":  Unique,
	"7":  Community,
	"8":  Valve,
	"9":  SelfMade,
	"10": Customized,
	"11": Strange,
	"12": Completed,
	"13": Haunted,
	"14": Collectors,
	"15": Decorated,
}

type PricingData struct {
	Response struct {
		Success          int64          `json:"success,omitempty"`
		CurrentTime      int64          `json:"current_time,omitempty"`
		RawUsdValue      float64        `json:"raw_usd_value,omitempty"`
		UsdCurrency      string         `json:"usd_currency,omitempty"`
		UsdCurrencyIndex int64          `json:"usd_currency_index,omitempty"`
		Items            map[string]any `json:"items,omitempty"`
	} `json:"response,omitempty"`
}

type PricingDataCache struct {
	CachedOn time.Time `json:"cachedOn"`
	Items    map[string]struct {
		Prices map[QualityItems]struct {
			Craftable   map[int]Item `json:"craftable"`
			Uncraftable map[int]Item `json:"non-craftable"`
		} `json:"prices"`
	} `json:"items"`
}

type Item struct {
	Value    float64      `json:"value"`
	Currency PricingItems `json:"currency"`
}

func (pd *PricingData) toCache() (*PricingDataCache, error) {
	cache := &PricingDataCache{
		CachedOn: time.Now(),
		Items: make(map[string]struct {
			Prices map[QualityItems]struct {
				Craftable   map[int]Item `json:"craftable"`
				Uncraftable map[int]Item `json:"non-craftable"`
			} `json:"prices"`
		}),
	}

	for itemName, itemData := range pd.Response.Items {
		itemMap, ok := itemData.(map[string]any)
		if !ok {
			continue
		}

		// Initialize item in cache
		cacheItem := struct {
			Prices map[QualityItems]struct {
				Craftable   map[int]Item `json:"craftable"`
				Uncraftable map[int]Item `json:"non-craftable"`
			} `json:"prices"`
		}{
			Prices: make(map[QualityItems]struct {
				Craftable   map[int]Item `json:"craftable"`
				Uncraftable map[int]Item `json:"non-craftable"`
			}),
		}

		// Process prices
		if prices, ok := itemMap["prices"].(map[string]any); ok {
			for qualityStr, qualityData := range prices {
				quality, ok := qualityMap[qualityStr]
				if !ok {
					continue
				}

				qualityMapData, ok := qualityData.(map[string]any)
				if !ok {
					continue
				}

				// Initialize quality entry
				qualityEntry := struct {
					Craftable   map[int]Item `json:"craftable"`
					Uncraftable map[int]Item `json:"non-craftable"`
				}{
					Craftable:   make(map[int]Item),
					Uncraftable: make(map[int]Item),
				}

				// Process Tradable and Non-Tradable
				if tradableData, ok := qualityMapData["Tradable"].(map[string]any); ok {
					for _, craftableKey := range []string{"Craftable", "Non-craftable"} {
						if _, keyExists := tradableData[craftableKey]; !keyExists {
							continue
						}
						if craftableData, ok := tradableData[craftableKey].(map[string]any); ok {
							for key, valueData := range craftableData {
								keyInt, err := strconv.Atoi(key)
								if err != nil {
									keyInt = 0
								}

								valueMap, ok := valueData.(map[string]any)
								if !ok {
									continue
								}

								item := Item{
									Currency: PricingItems(valueMap["currency"].(string)),
									Value:    valueMap["value"].(float64),
								}

								if craftableKey == "Craftable" {
									qualityEntry.Craftable[keyInt] = item
								} else {
									qualityEntry.Uncraftable[keyInt] = item
								}
							}
						} else if singleCraftableData, ok := tradableData[craftableKey].([]any); ok {
							valueMap, exists := singleCraftableData[0].(map[string]any)
							if !exists {
								continue
							}
							item := Item{
								Currency: PricingItems(valueMap["currency"].(string)),
								Value:    valueMap["value"].(float64),
							}

							if craftableKey == "Craftable" {
								qualityEntry.Craftable[0] = item
							} else {
								qualityEntry.Uncraftable[0] = item
							}
						} else {
							log.Printf("Unprocessable: %s: %s", craftableKey, itemName)
						}
					}
				}
				cacheItem.Prices[quality] = qualityEntry
			}
		}
		cache.Items[itemName] = cacheItem
	}

	return cache, nil
}

type CurrencyData struct {
	Response struct {
		Success    int64 `json:"success,omitempty"`
		Currencies map[PricingItems]struct {
			Name  string `json:"name"`
			Price struct {
				Value    float64 `json:"value"`
				Currency string  `json:"currency"`
			}
		} `json:"currencies,omitempty"`
	} `json:"response,omitempty"`
}
