package backpack

const minimumPriceInCents = 5

const (
	Metal = "metal"
	Keys  = "keys"
	Usd   = "usd"
)

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

type flatCurrency struct {
	Keys  float64
	Metal float64
}

func (c *currencyData) flatten() *flatCurrency {
	metalValue := c.Response.Currencies["metal"].Price.Value
	keyValue := c.Response.Currencies["keys"].Price.Value

	return &flatCurrency{
		Metal: metalValue,
		Keys:  keyValue * metalValue,
	}
}

func (c *flatCurrency) toRealPrice(value float64, currency string) uint {
	var price uint
	switch currency {
	case Keys:
		price = uint(c.Keys * value * 100)
	case Metal:
		price = uint(c.Metal * value * 100)
	case Usd:
		price = uint(value * 100)
	}

	return max(price, minimumPriceInCents)
}
