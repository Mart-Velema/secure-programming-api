package inventory

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"guineatrade.nhlstenden.com/src/auth/middleware"
	"guineatrade.nhlstenden.com/src/steam"
)

func GetUserStock(c *gin.Context) {
	user, err := middleware.ExtractTokenUser(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Token expired"})
		return
	}

	userInventory, err := GetUserInventory(user.SteamId)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Could not get user inventory"})
		return
	}

	c.JSON(http.StatusOK, userInventory.ToItem().ToStock())
}

func GetSteamBotStock(c *gin.Context) {
	botInventory, err := steam.GetBotInventoryData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get bot inventory"})
		return
	}

	c.JSON(http.StatusOK, botInventory.ToStock())
}
