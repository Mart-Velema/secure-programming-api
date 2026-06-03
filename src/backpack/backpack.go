package backpack

import "github.com/gin-gonic/gin"

func GetPrices(c *gin.Context) {
	if PricingCache.CachedOn.IsZero() {
		c.String(503, "Cache not initialised")
		return
	}
	c.JSON(200, PricingCache)
}
