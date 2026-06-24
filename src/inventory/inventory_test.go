package inventory

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	_ "guineatrade.nhlstenden.com/src/1nit"
	"guineatrade.nhlstenden.com/src/auth"
)

// Steam ID for testing: 76561198248575244
// This ID belongs to a MPTF bot, and will always be available
func TestGetInventoryApi(t *testing.T) {
	router := gin.Default()
	router.POST("/api/v1/auth/register", auth.Register)
	router.POST("/api/v1/auth/login", auth.Login)
	router.PATCH("/api/v1/auth/steam", auth.UpdateSteam)
	router.GET("/api/v1/auth/me", auth.Me)
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

	password := uuid.New()
	email := fmt.Sprintf("%s@%s.com", uuid.New(), uuid.New())

	user := struct {
		Email          string `json:"email"`
		Name           string `json:"name"`
		Password       string `json:"password"`
		PasswordVerify string `json:"passwordVerify"`
	}{
		Email:          email,
		Name:           rand.Text(),
		Password:       password.String(),
		PasswordVerify: password.String(),
	}

	userJson, _ := json.Marshal(&user)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(userJson))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(userJson))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var tokens = auth.Tokens{}
	err := json.Unmarshal(w.Body.Bytes(), &tokens)
	assert.IsEqual(err, nil)

	for _, test := range tests {
		tradeUrl := fmt.Sprintf("https://steampowered.com/user/%d/trade", test.Id)
		patchSteam := struct {
			TradeUrl string `json:"tradeUrl"`
			SteamId  uint64 `json:"steamId"`
		}{
			TradeUrl: tradeUrl,
			SteamId:  test.Id,
		}

		w = httptest.NewRecorder()
		patchSteamJson, _ := json.Marshal(patchSteam)
		req, _ = http.NewRequest("PATCH", "/api/v1/auth/steam", bytes.NewReader(patchSteamJson))
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens.JWT))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/v1/user/inventory", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens.JWT))
		router.ServeHTTP(w, req)

		if test.HasError {
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			assert.Equal(t, "{\"error\":\"Unable to get inventory data\"}", w.Body.String())
			continue
		}
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, w.Body.Len() > 1, true)
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/user/inventory", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", "not-a-valid-token"))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, "{\"error\":\"Token expired\"}", w.Body.String())
}
