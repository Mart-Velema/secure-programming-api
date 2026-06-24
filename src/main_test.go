package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	mathrand "math/rand/v2"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"guineatrade.nhlstenden.com/src/auth"
	"guineatrade.nhlstenden.com/src/auth/mfa"
)

func TestAuthFlow(t *testing.T) {
	router := CreateRouter()

	// ===========
	// Register
	w := httptest.NewRecorder()
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
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(userJson))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// ===========
	// Login
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(userJson))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var tokens = auth.Tokens{}
	err := json.Unmarshal(w.Body.Bytes(), &tokens)
	assert.IsEqual(err, nil)

	w = httptest.NewRecorder()
	tokenJson, _ := json.Marshal(tokens)
	req, _ = http.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewReader(tokenJson))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// ===========
	// Get user profile
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/auth/me", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens.JWT))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var meUser struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Balance  int64  `json:"balance"`
		HasMFA   bool   `json:"mfaEnabled"`
		SteamId  uint64 `json:"steamId"`
		TradeUrl string `json:"tradeUrl"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &meUser)

	assert.Equal(t, meUser.Name, user.Name)
	assert.Equal(t, meUser.Email, user.Email)
	assert.Equal(t, meUser.HasMFA, false)
	assert.Equal(t, meUser.Balance, int64(0))
	assert.Equal(t, meUser.SteamId, uint64(0))
	assert.Equal(t, meUser.TradeUrl, "")

	user.Password = rand.Text()
	userJson, _ = json.Marshal(&user)

	// ===========
	// Update password

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/api/v1/auth/me", bytes.NewReader(userJson))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens.JWT))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)

	// ===========
	// Register Steam
	steamId := mathrand.Uint64() >> 1
	tradeUrl := fmt.Sprintf("https://steampowered.com/user/%d/trade", steamId)
	patchSteam := struct {
		TradeUrl string `json:"tradeUrl"`
		SteamId  uint64 `json:"steamId"`
	}{
		TradeUrl: tradeUrl,
		SteamId:  steamId,
	}

	w = httptest.NewRecorder()
	patchSteamJson, _ := json.Marshal(patchSteam)
	req, _ = http.NewRequest("PATCH", "/api/v1/auth/steam", bytes.NewReader(patchSteamJson))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens.JWT))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/auth/me", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens.JWT))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	_ = json.Unmarshal(w.Body.Bytes(), &meUser)

	assert.Equal(t, meUser.Name, user.Name)
	assert.Equal(t, meUser.Email, user.Email)
	assert.Equal(t, meUser.HasMFA, false)
	assert.Equal(t, meUser.Balance, int64(0))
	assert.Equal(t, meUser.SteamId, steamId)
	assert.Equal(t, meUser.TradeUrl, tradeUrl)

	// ===========
	// Register TOTP
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/mfa/totp/register", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens.JWT))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var totpCredentials mfa.TotpCodes
	_ = json.Unmarshal(w.Body.Bytes(), &totpCredentials)

	// ===========
	// Verify TOTP
	code, _ := totp.GenerateCode(totpCredentials.Code, time.Now())
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/mfa/totp/verify", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens.JWT))
	req.Header.Add("X-TOTP-Code", code)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// ===========
	// Reset TOTP TOTP Code
	code, _ = totp.GenerateCode(totpCredentials.Code, time.Now())
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/auth/mfa/totp/reset", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens.JWT))
	req.Header.Add("X-TOTP-Code", code)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// ===========
	// Register TOTP Again
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/mfa/totp/register", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens.JWT))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	_ = json.Unmarshal(w.Body.Bytes(), &totpCredentials)

	// ===========
	// Reset TOTP Recovery code
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/auth/mfa/totp/reset", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens.JWT))
	req.Header.Add("X-Recovery-Code", totpCredentials.RecoveryCode)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
