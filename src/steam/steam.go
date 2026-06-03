package steam

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func GetBotStatus(c *gin.Context) {
	botURL := os.Getenv("STEAM_BOT_URL")
	botAPIKey := os.Getenv("BOT_API_KEY")

	req, err := http.NewRequest("GET", botURL+"/steam/status", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Header.Set("X-API-Key", botAPIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(resp.StatusCode, "application/json", body)
}