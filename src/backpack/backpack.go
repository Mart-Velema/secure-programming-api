package backpack

import "github.com/gin-gonic/gin"

func GetPrices(c *gin.Context) {
	if PricingCache.CachedOn.IsZero() {
		c.String(404, "Cache not initialised")
	}
	c.JSON(200, PricingCache)
}

func GetItemDetails(c *gin.Context) {
	itemId := c.Param("item")
	c.String(200, itemId)
}

func GetCurrencies(c *gin.Context) {
	if CurrencyCache.CachedOn.IsZero() {
		c.String(404, "Cache not initialised")
	}
	c.JSON(200, CurrencyCache)
}
