package inventory

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
	router.GET("/api/v1/user/inventory", GetInventory)

	type getInventoryTest struct {
		Id                uint64
		InventoryResponse int
		HasError          bool
	}
	tests := []getInventoryTest{
		{76561198248575244, 1, false},
		{0, 0, true},
	}

	user := database.CreateRandomUser()
	database.GetInstance().First(user)
	jwt, _ := middleware.GenerateToken(user)

	for _, test := range tests {
		user.SteamId = test.Id
		user.TradeUrl = fmt.Sprintf("https://steampowered.com/user/%d/trade", test.Id)
		database.GetInstance().Select("steam_id", "trade_url").Save(user)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/user/inventory", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
		router.ServeHTTP(w, req)

		if test.HasError {
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			assert.Equal(t, "{\"error\":\"Unable to get inventory data\"}", w.Body.String())
			continue
		}
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, w.Body.Len() > 1, true)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/user/inventory", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", "not-a-valid-token"))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, "{\"error\":\"Token expired\"}", w.Body.String())
}
