package inventory

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	_ "guineatrade.nhlstenden.com/src/1nit"
	"guineatrade.nhlstenden.com/src/auth/middleware"
	"guineatrade.nhlstenden.com/src/database"
)

// Steam ID for testing: 76561198248575244
// This ID belongs to a MPTF bot, and will always be available
func TestGetInventoryApi(t *testing.T) {
	router := gin.Default()
	router.GET("/inventory", GetInventory)

	type getInventoryTest struct {
		Id         uint64
		StatusCode int
		Response   string
	}
	tests := []getInventoryTest{
		{76561198248575244, http.StatusOK, ""},
		{0, http.StatusInternalServerError, "{\"error\":\"Unable to get inventory data\"}"},
	}

	user := database.CreateRandomUser()
	database.GetInstance().First(user)
	jwt, _ := middleware.GenerateToken(user)

	for _, test := range tests {
		user.SteamId = test.Id
		user.TradeUrl = fmt.Sprintf("https://steampowered.com/user/%d/trade", test.Id)
		database.GetInstance().Select("steam_id", "trade_url").Save(user)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/inventory", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
		router.ServeHTTP(w, req)

		assert.Equal(t, strings.Contains(w.Body.String(), test.Response), true)
		assert.Equal(t, test.StatusCode, w.Code)
	}
}
