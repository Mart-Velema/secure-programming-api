package backpack

import "github.com/gin-gonic/gin"

func GetPrices(c *gin.Context) {
	if PricingCache.CachedOn.IsZero() {
		c.String(503, "Cache not initialised")
		return
	}
	c.JSON(200, PricingCache)
}

func GetItemDetails(c *gin.Context) {
	if PricingCache.CachedOn.IsZero() {
		c.String(503, "Cache not initialised")
		return
	}
	itemId := c.Param("item")
	itemPricing, ok := PricingCache.Items[itemId]
	if !ok {
		c.String(404, "Item not found")
		return
	}
	c.JSON(200, itemPricing)
}

func GetCurrencies(c *gin.Context) {
	if CurrencyCache.CachedOn.IsZero() {
		c.String(503, "Cache not initialised")
		return
	}
	c.JSON(200, CurrencyCache)
}
