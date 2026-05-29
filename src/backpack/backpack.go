package backpack

import "github.com/gin-gonic/gin"

func GetPrices(c *gin.Context) {
	c.Status(200)
}

func GetItemDetails(c *gin.Context) {
	itemId := c.Param("item")
	c.String(200, itemId)
}

func GetCurrencies(c *gin.Context) {
	c.Status(200)
}
