package backpack

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

type Quality string

const (
	Normal     Quality = "Normal"
	Genuine    Quality = "Genuine"
	Rarity2    Quality = "Rarity2"
	Vintage    Quality = "Vintage"
	Rarity3    Quality = "Rarity3"
	Unusual    Quality = "Unusual"
	Unique     Quality = "Unique"
	Community  Quality = "Community"
	Valve      Quality = "Valve"
	SelfMade   Quality = "Self-Made"
	Customized Quality = "Customized"
	Strange    Quality = "Strange"
	Completed  Quality = "Completed"
	Haunted    Quality = "Haunted"
	Collectors Quality = "Collectors"
	Decorated  Quality = "Decorated"
)

var qualityMap = map[string]Quality{
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

type multiplier uint

const (
	None           multiplier = 1
	ScrapMetal     multiplier = 20
	RefinedMetal   multiplier = 20
	ReclaimedMetal multiplier = 20
)

// API output
type pricingData struct {
	Response struct {
		Success          int64               `json:"success,omitempty"`
		CurrentTime      int64               `json:"current_time,omitempty"`
		RawUsdValue      float64             `json:"raw_usd_value,omitempty"`
		UsdCurrency      string              `json:"usd_currency,omitempty"`
		UsdCurrencyIndex int64               `json:"usd_currency_index,omitempty"`
		Items            map[string]itemData `json:"items,omitempty"`
	} `json:"response,omitempty"`
}

func (pd *pricingData) toCache(currencyConversions *flatCurrency) (*PricingDataCache, error) {
	cache := &PricingDataCache{
		CachedOn: time.Now(),
		Items:    make(map[string]ItemDetails),
	}

	for itemName, item := range pd.Response.Items {
		if len(item.Defindex) == 0 {
			log.Printf("Can't find defindexes for: %s", itemName)
			continue
		}

		var batchSize multiplier
		switch itemName {
		case "Scrap Metal":
			batchSize = ScrapMetal
		case "Refined Metal":
			batchSize = RefinedMetal
		case "Reclaimed Metal":
			batchSize = ReclaimedMetal
		default:
			batchSize = None
		}

		var defindexList = make([]uint, len(item.Defindex))
		for idx, defindex := range item.Defindex {
			defindexList[idx] = uint(defindex)
		}
		itemConstant, ok := itemCache[strconv.Itoa(int(defindexList[0]))] // Because of course just using uint is not good enough for you
		if !ok {
			log.Printf("Can't decode defindex: %s", itemName)
			continue
		}

		cacheItem := ItemDetails{
			IconUrl:        itemConstant.Url,
			Defindex:       defindexList,
			MarketHashName: itemConstant.MarketHashName,
			Prices:         make(map[Quality]ItemPair),
		}
		err := cacheItem.toCache(item.Prices, currencyConversions, batchSize)
		if err != nil {
			log.Printf("Unable to cache item: %s: %s", itemName, err)
			continue
		}

		cache.Items[itemName] = cacheItem
	}

	fmt.Println(cache.Items["Refined Metal"].Prices[Unique].Craftable)
	fmt.Println(cache.Items["Reclaimed Metal"].Prices[Unique].Craftable)
	fmt.Println(cache.Items["Scrap Metal"].Prices[Unique].Craftable)

	return cache, nil
}

type itemData struct {
	Defindex []float64           `json:"defindex"`
	Prices   map[string]tradable `json:"prices"`
}

type tradable struct {
	Tradable struct {
		Craftable   any `json:"Craftable"`
		Uncraftable any `json:"Non-Craftable"`
	} `json:"Tradable"` // EVIL!!!
}

func parseField(field any, currencyConversions *flatCurrency, batchSize multiplier) (map[uint]uint, error) {
	resultMap := make(map[uint]uint)

	switch v := field.(type) {
	case map[string]any:
		parsed, err := parseUnusual(v)
		if err != nil {
			return nil, err
		}
		for idx, rawPrice := range parsed {
			rawPrice.Value *= float64(batchSize)
			resultMap[idx] = currencyConversions.toRealPrice(rawPrice.Value, rawPrice.Currency)
		}

	case []any:
		if len(v) == 0 {
			return nil, errors.New("item field empty")
		}
		result, err := parseOther(v[0])
		if err != nil {
			return nil, err
		}

		result.Value *= float64(batchSize)
		resultMap[0] = currencyConversions.toRealPrice(result.Value, result.Currency)
	}

	return resultMap, nil
}

type rawPrices struct {
	Value    float64
	Currency string
}

func parseOther(rawItemDetails any) (rawPrices, error) {
	var rawPrice rawPrices
	valueMap, ok := rawItemDetails.(map[string]any)
	if !ok {
		return rawPrice, errors.New("input data is not a valid item")
	}

	if valueMap["currency"] == nil || valueMap["value"] == nil {
		return rawPrice, errors.New("input data does not contain value or currency fields")
	}
	rawPrice.Value = valueMap["value"].(float64)
	rawPrice.Currency = valueMap["currency"].(string)

	return rawPrice, nil
}

func parseUnusual(rawItemDetails map[string]any) (map[uint]rawPrices, error) {
	rawPrice := make(map[uint]rawPrices)

	for key, valueData := range rawItemDetails {
		keyInt, err := strconv.Atoi(key)
		if err != nil {
			keyInt = 0
		}
		singlePrice, err := parseOther(valueData)
		if err != nil {
			return nil, err
		}
		rawPrice[uint(keyInt)] = singlePrice
	}

	return rawPrice, nil
}

type PricingDataCache struct {
	CachedOn time.Time              `json:"cachedOn"`
	Items    map[string]ItemDetails `json:"items"`
}

type ItemDetails struct {
	IconUrl        string               `json:"icon"`
	Defindex       []uint               `json:"defindex"`
	MarketHashName string               `json:"marketHashName"`
	Prices         map[Quality]ItemPair `json:"prices"`
}

func (id *ItemDetails) toCache(tradable map[string]tradable, currencyConversion *flatCurrency, batchSize multiplier) error {
	for qualityNumber, tradableItem := range tradable {
		qualityString, ok := qualityMap[qualityNumber]
		if !ok {
			return fmt.Errorf("unknown quality number: %s", qualityNumber)
		}
		itemPair := ItemPair{
			Craftable:   make(map[uint]uint),
			Uncraftable: make(map[uint]uint),
		}

		// Total hack, will only work if we keep multiplier only for unique items.
		// But then again, who would want to buy 20 vintage metal for almost $12k...
		if qualityString != Unique {
			batchSize = None
		}

		err := itemPair.toCache(tradableItem, currencyConversion, batchSize)
		if err != nil {
			return err
		}
		id.Prices[qualityString] = itemPair
	}

	return nil
}

type ItemPair struct {
	Craftable   map[uint]uint `json:"craftable,omitempty"`
	Uncraftable map[uint]uint `json:"non-craftable,omitempty"`
}

func (ip *ItemPair) toCache(t tradable, currencyConversion *flatCurrency, batchSize multiplier) error {
	craftableMap, err := parseField(t.Tradable.Craftable, currencyConversion, batchSize)
	if err != nil {
		return err
	}
	ip.Craftable = craftableMap

	uncraftableMap, err := parseField(t.Tradable.Uncraftable, currencyConversion, batchSize)
	if err != nil {
		return err
	}
	ip.Uncraftable = uncraftableMap

	return nil
}
