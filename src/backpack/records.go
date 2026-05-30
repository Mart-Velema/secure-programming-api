package backpack

type PricingItems string

const (
	PricingMetal   = "metal"
	PricingKey     = "keys"
	PricingEarbuds = "earbuds"
)

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
