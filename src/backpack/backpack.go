package backpack

import "github.com/gin-gonic/gin"

func GetPrices(c *gin.Context) {
	if PricingCache.CachedOn.IsZero() {
		c.String(503, "Cache not initialised")
		return
	}
	c.JSON(200, PricingCache)
}

func GetSpecificPrice(marketHashName string, quality Quality, isCraftable bool, effectId string) uint {
	if PricingCache.CachedOn.IsZero() {
		return 0
	}

	itemPair := PricingCache.Items[marketHashName].Prices[quality]

	if isCraftable {
		return itemPair.Craftable[effectId]
	}
	return itemPair.Uncraftable[effectId]
}
