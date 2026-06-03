package backpack

import (
	"errors"
	"log"
	"strconv"
	"time"
)

type qualityItems int

const (
	Normal qualityItems = iota
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

var qualityMap = map[string]qualityItems{
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

type pricingData struct {
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
	CachedOn time.Time              `json:"cachedOn"`
	Items    map[string]ItemDetails `json:"items"`
}

type ItemDetails struct {
	IconUrl  string                    `json:"icon"`
	Defindex []uint                    `json:"defindex"`
	Prices   map[qualityItems]ItemPair `json:"prices"`
}

type ItemPair struct {
	Craftable   map[int]uint `json:"craftable,omitempty"`
	Uncraftable map[int]uint `json:"non-craftable,omitempty"`
}

func (ip *ItemPair) addItem(key int, valueData any, isCraftable bool, currencyConversions map[string]float64) error {
	valueMap, ok := valueData.(map[string]any)
	if !ok {
		return errors.New("input data is not a valid item")
	}

	if valueMap["currency"] == nil || valueMap["value"] == nil {
		return errors.New("input data does not contain value or currency fields")
	}

	price := currencyConversions[valueMap["currency"].(string)] * valueMap["value"].(float64)
	price *= 100

	if isCraftable {
		ip.Craftable[key] = uint(price)
	} else {
		ip.Uncraftable[key] = uint(price)
	}

	return nil
}

func (pd *pricingData) toCache(currencyConversions map[string]float64) (*PricingDataCache, error) {
	cache := &PricingDataCache{
		CachedOn: time.Now(),
		Items:    make(map[string]ItemDetails),
	}

	for itemName, itemData := range pd.Response.Items {
		itemMap, ok := itemData.(map[string]any)
		if !ok {
			continue
		}
		prices, ok := itemMap["prices"].(map[string]any)
		if !ok {
			continue
		}

		defindexes, ok := itemMap["defindex"].([]any)
		if !ok || len(defindexes) == 0 {
			log.Printf("Can't find defindex: %s", itemName)
			continue
		}
		var defindexList = make([]uint, len(defindexes))
		for idx, defindex := range defindexes {
			defindexList[idx] = uint(defindex.(float64))
		}

		itemUrl, ok := itemCache[strconv.Itoa(int(defindexList[0]))] // Because of course just using uint is not good enough for you
		if !ok {
			log.Printf("Can't decode defindex: %s", itemName)
			continue
		}

		cacheItem := ItemDetails{
			IconUrl:  itemUrl,
			Defindex: defindexList,
			Prices:   make(map[qualityItems]ItemPair),
		}
		for qualityStr, qualityData := range prices {
			quality, ok := qualityMap[qualityStr]
			if !ok {
				continue
			}

			qualityMapData, ok := qualityData.(map[string]any)
			if !ok {
				continue
			}

			itemPair, err := createItemPair(qualityMapData, currencyConversions)
			if err != nil {
				log.Printf("%s: %s", err, itemName)
				continue
			}

			cacheItem.Prices[quality] = itemPair
		}
		if len(cacheItem.Prices) != 0 {
			cache.Items[itemName] = cacheItem
		}
	}

	return cache, nil
}

func createItemPair(qualityMapData map[string]any, currencyConversions map[string]float64) (ItemPair, error) {
	itemPair := ItemPair{
		Craftable:   make(map[int]uint),
		Uncraftable: make(map[int]uint),
	}

	tradableData, ok := qualityMapData["Tradable"].(map[string]any)
	if !ok {
		return itemPair, errors.New("no Tradable key found")
	}

	for _, craftableKey := range []string{"Craftable", "Non-Craftable"} {
		if _, keyExists := tradableData[craftableKey]; !keyExists {
			continue
		}
		if craftableData, ok := tradableData[craftableKey].(map[string]any); ok {
			for key, valueData := range craftableData {
				keyInt, err := strconv.Atoi(key)
				if err != nil {
					keyInt = 0
				}
				err = itemPair.addItem(keyInt, valueData, craftableKey == "Craftable", currencyConversions)
				if err != nil {
					return ItemPair{}, err
				}
			}
		} else if singleCraftableData, ok := tradableData[craftableKey].([]any); ok {
			err := itemPair.addItem(0, singleCraftableData[0], craftableKey == "Craftable", currencyConversions)
			if err != nil {
				return ItemPair{}, err
			}
		} else {
			return itemPair, errors.New("not a valid item")
		}
	}

	return itemPair, nil
}

type currencyData struct {
	Response struct {
		Success    int64 `json:"success,omitempty"`
		Currencies map[string]struct {
			Name  string `json:"name"`
			Price struct {
				Value    float64 `json:"value"`
				Currency string  `json:"currency"`
			}
		} `json:"currencies,omitempty"`
	} `json:"response,omitempty"`
}

func (c *currencyData) flatten() map[string]float64 {
	currencyCache := make(map[string]float64)

	metalValue := c.Response.Currencies["metal"].Price.Value
	keyValue := c.Response.Currencies["keys"].Price.Value

	currencyCache["metal"] = metalValue
	currencyCache["keys"] = keyValue * metalValue

	return currencyCache
}
