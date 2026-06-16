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
