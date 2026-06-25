package steam

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"guineatrade.nhlstenden.com/src/items"
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
	result, err := GetBotInventoryData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get inventory data"})
		return
	}
	c.JSON(http.StatusOK, result.Assets)
}

func GetBotInventoryData() (items.Items, error) {
	const appId = 440
	const contextId = 2

	path := fmt.Sprintf(
		"/steam/inventory?appId=%d&contextId=%d",
		appId,
		contextId,
	)

	result, err := steamBotRequest(http.MethodGet, path, nil)
	if err != nil {
		return items.Items{}, err
	}
	var inventory items.SteamInventoryResponse
	err = json.Unmarshal(*result, &inventory)
	if err != nil {
		return items.Items{}, err
	}
	return inventory.ToItems(), nil
}

func GetTradeOffers(c *gin.Context) {
	_, err := steamBotRequest(http.MethodGet, "/steam/trade-offers", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get trade offers"})
		return
	}
}

func GetTradeOfferHistory(c *gin.Context) {
	_, err := steamBotRequest(http.MethodGet, "/steam/trade-offers/history", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get trade offer history"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get requested trade offer"})
		return
	}
}

func SendTradeOffer(tradeOfferRequest SendTradeOfferRequest) (*SendTradeOfferResponse, error) {
	if tradeOfferRequest.TradeURL == "" {
		return nil, errors.New("tradeUrl is required")
	}

	if len(tradeOfferRequest.ItemsToGive) == 0 && len(tradeOfferRequest.ItemsToReceive) == 0 {
		return nil, errors.New("trade offer must contain at least one item")
	}

	body, err := json.Marshal(tradeOfferRequest)
	if err != nil {
		return nil, err
	}

	result, err := steamBotRequest(http.MethodPost, "/steam/trade-offers", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var response SendTradeOfferResponse
	if err = json.Unmarshal(*result, &response); err != nil {
		return nil, err
	}

	if !response.OK {
		return nil, errors.New("unable to process trade offer request")
	}

	return &response, nil
}

func AcceptTradeOffer(c *gin.Context) {
	tradeOfferId := c.Param("tradeOfferId")

	path := fmt.Sprintf(
		"/steam/trade-offers/%s/accept",
		tradeOfferId,
	)

	_, err := steamBotRequest(http.MethodPost, path, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to accept trade offers"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to cancel trade offers"})
		return
	}
}
