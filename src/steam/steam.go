package steam

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"log"

	"github.com/gin-gonic/gin"
)

func steamBotRequest(c *gin.Context, method string, path string, body io.Reader) {
	botURL := os.Getenv("STEAM_BOT_URL")
	botAPIKey := os.Getenv("BOT_API_KEY")

	url := fmt.Sprintf("%s%s", botURL, path)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Header.Set("X-API-Key", botAPIKey)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer func(Body io.ReadCloser) {
        err := Body.Close()
        if err != nil {
            log.Println(err)
        }
    }(resp.Body)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(resp.StatusCode, "application/json", responseBody)
}

func GetBotStatus(c *gin.Context) {
	steamBotRequest(c, http.MethodGet, "/steam/status", nil)
}

func GetBotInventory(c *gin.Context) {
	appId := c.DefaultQuery("appId", "730")
	contextId := c.DefaultQuery("contextId", "2")

	path := fmt.Sprintf(
		"/steam/inventory?appId=%s&contextId=%s",
		appId,
		contextId,
	)

	steamBotRequest(c, http.MethodGet, path, nil)
}

func GetTradeOffers(c *gin.Context) {
	steamBotRequest(c, http.MethodGet, "/steam/trade-offers", nil)
}

func GetTradeOfferHistory(c *gin.Context) {
	steamBotRequest(c, http.MethodGet, "/steam/trade-offers/history", nil)
}

func GetTradeOffer(c *gin.Context) {
	tradeOfferId := c.Param("tradeOfferId")

	path := fmt.Sprintf(
		"/steam/trade-offers/%s",
		tradeOfferId,
	)

	steamBotRequest(c, http.MethodGet, path, nil)
}

func LoginBot(c *gin.Context) {
	steamBotRequest(c, http.MethodPost, "/steam/login", c.Request.Body)
}

func SendTradeOffer(c *gin.Context) {
	steamBotRequest(c, http.MethodPost, "/steam/trade-offers", c.Request.Body)
}

func AcceptTradeOffer(c *gin.Context) {
	tradeOfferId := c.Param("tradeOfferId")

	path := fmt.Sprintf(
		"/steam/trade-offers/%s/accept",
		tradeOfferId,
	)

	steamBotRequest(c, http.MethodPost, path, nil)
}

func CancelTradeOffer(c *gin.Context) {
	tradeOfferId := c.Param("tradeOfferId")

	path := fmt.Sprintf(
		"/steam/trade-offers/%s/cancel",
		tradeOfferId,
	)

	steamBotRequest(c, http.MethodPost, path, nil)
}