package steam

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func steamBotRequest(method string, path string, body io.Reader) (*[]byte, error) {
	botURL := os.Getenv("STEAM_BOT_URL")
	botAPIKey := os.Getenv("BOT_API_KEY")

	url := fmt.Sprintf("%s%s", botURL, path)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", botAPIKey)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

func GetBotStatus(c *gin.Context) {
	_, err := steamBotRequest(http.MethodGet, "/steam/status", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func GetBotInventory(c *gin.Context) {
	appId := c.DefaultQuery("appId", "440")
	contextId := c.DefaultQuery("contextId", "2")

	if appId != "440" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unsupported appId",
		})
		return
	}

	if contextId != "2" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unsupported contextId",
		})
		return
	}

	path := fmt.Sprintf(
		"/steam/inventory?appId=%s&contextId=%s",
		appId,
		contextId,
	)

	// TODO: Parse properly
	_, err := steamBotRequest(http.MethodGet, path, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func GetTradeOffers(c *gin.Context) {
	_, err := steamBotRequest(http.MethodGet, "/steam/trade-offers", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func GetTradeOfferHistory(c *gin.Context) {
	_, err := steamBotRequest(http.MethodGet, "/steam/trade-offers/history", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func GetTradeOffer(c *gin.Context) {
	tradeOfferId := c.Param("tradeOfferId")

	path := fmt.Sprintf(
		"/steam/trade-offers/%s",
		tradeOfferId,
	)

	_, err := steamBotRequest(http.MethodGet, path, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func SendTradeOffer(c *gin.Context) {
	_, err := steamBotRequest(http.MethodPost, "/steam/trade-offers", c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func AcceptTradeOffer(c *gin.Context) {
	tradeOfferId := c.Param("tradeOfferId")

	path := fmt.Sprintf(
		"/steam/trade-offers/%s/accept",
		tradeOfferId,
	)

	_, err := steamBotRequest(http.MethodPost, path, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func CancelTradeOffer(c *gin.Context) {
	tradeOfferId := c.Param("tradeOfferId")

	path := fmt.Sprintf(
		"/steam/trade-offers/%s/cancel",
		tradeOfferId,
	)

	_, err := steamBotRequest(http.MethodPost, path, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
