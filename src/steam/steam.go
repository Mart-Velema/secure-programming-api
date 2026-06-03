package steam

import (
	"fmt"
	"io"
	"net/http"
	"os"

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
	defer resp.Body.Close()

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