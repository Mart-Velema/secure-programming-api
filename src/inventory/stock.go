package inventory

import (
	"log"
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

	log.Println(user.SteamId)
	userInventory, err := getInventory(user.SteamId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Could not get user inventory"})
		return
	}

	c.JSON(http.StatusOK, userInventory.ToItem().ToStock())
}

func GetSteamBotStock(c *gin.Context) {
	botInventory, err := steam.GetBotInventoryData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get bot inventory"})
	}

	c.JSON(http.StatusOK, botInventory.ToStock())
}
