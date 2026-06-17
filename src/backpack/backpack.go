package backpack

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPrices(c *gin.Context) {
	if PricingCache.CachedOn.IsZero() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Cache not instantiated"})
		return
	}
	c.JSON(http.StatusOK, PricingCache)
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
