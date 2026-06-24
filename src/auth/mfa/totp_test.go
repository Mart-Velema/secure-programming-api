package mfa

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/pquerna/otp/totp"
	_ "guineatrade.nhlstenden.com/src/1nit"
	"guineatrade.nhlstenden.com/src/auth/middleware"
	"guineatrade.nhlstenden.com/src/database"
)

func TestRegisterTOTP(t *testing.T) {
	router := gin.Default()
	router.POST("/register", RegisterTOTP)

	user := database.CreateRandomUser()
	jwt, _ := middleware.GenerateToken(user)

	type test struct {
		JWT        string
		StatusCode int
		Result     string
	}

	tests := []test{
		{jwt, http.StatusOK, ""},
		{jwt, http.StatusConflict, "{\"error\":\"TOTP already registered\"}"},
		{"", http.StatusUnprocessableEntity, "{\"error\":\"Token expired\"}"},
	}

	user.RecoveryCode = ""
	user.TotpSecret = ""
	database.GetInstance().Select("totp_secret", "recovery_code").Save(user)

	for _, t2 := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t2.JWT))
		router.ServeHTTP(w, req)

		assert.Equal(t, t2.StatusCode, w.Code)
		if len(t2.Result) > 0 {
			assert.Equal(t, t2.Result, w.Body.String())
		}
	}
}

func TestVerifyTOTP(t *testing.T) {
	router := gin.Default()
	router.POST("/verify", VerifyTOTP)

	user := database.CreateRandomUser()
	jwt, _ := middleware.GenerateToken(user)

	type test struct {
		JWT        string
		TOTPCode   string
		StatusCode int
		Result     string
	}

	code, _ := totp.GenerateCode(user.TotpSecret, time.Now())
	tests := []test{
		{jwt, code, http.StatusOK, ""},
		{"", code, http.StatusUnprocessableEntity, "{\"error\":\"Token expired\"}"},
		{jwt, "", http.StatusUnauthorized, "{\"error\":\"No MFA supplied\"}"},
		{jwt, "123456", http.StatusUnauthorized, "{\"error\":\"Invalid MFA\"}"},
	}

	for _, t2 := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/verify", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t2.JWT))
		req.Header.Add("X-TOTP-Code", t2.TOTPCode)
		router.ServeHTTP(w, req)

		assert.Equal(t, t2.StatusCode, w.Code)
		if len(t2.Result) > 0 {
			assert.Equal(t, t2.Result, w.Body.String())
		}
	}
}

func TestResetTOTP(t *testing.T) {
	router := gin.Default()
	router.POST("/reset", ResetTOTP)

	user := database.CreateRandomUser()
	jwt, _ := middleware.GenerateToken(user)

	type test struct {
		JWT        string
		IsTOTPCode bool
		AddHeader  bool
		BadKey     bool
		StatusCode int
		Result     string
	}

	tests := []test{
		{jwt, true, true, false, http.StatusNoContent, ""},
		{"", true, true, false, http.StatusUnprocessableEntity, "{\"error\":\"Token expired\"}"},
		{jwt, true, true, true, http.StatusUnauthorized, "{\"error\":\"Invalid MFA\"}"},
		{jwt, false, true, false, http.StatusNoContent, ""},
		{jwt, false, false, false, http.StatusUnprocessableEntity, "{\"error\":\"No recovery or TOTP code supplied\"}"},
		{jwt, false, true, true, http.StatusUnauthorized, "{\"error\":\"Invalid recovery code\"}"},
	}

	for idx, t2 := range tests {
		fmt.Println(idx)
		user.TotpSecret = rand.Text()
		user.RecoveryCode = rand.Text()

		database.GetInstance().Select("recovery_code", "totp_secret").Save(user)
		database.GetInstance().First(user)

		var code string
		if t2.IsTOTPCode {
			code, _ = totp.GenerateCode(user.TotpSecret, time.Now())
		} else {
			code = user.RecoveryCode
		}
		if t2.BadKey {
			code = "abcdef"
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/reset", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t2.JWT))
		if t2.AddHeader {
			if t2.IsTOTPCode {
				req.Header.Add("X-TOTP-Code", code)
			} else {
				req.Header.Add("X-Recovery-Code", code)
			}
		}
		router.ServeHTTP(w, req)

		assert.Equal(t, t2.StatusCode, w.Code)
		if len(t2.Result) > 0 {
			assert.Equal(t, t2.Result, w.Body.String())
		}
	}
}
