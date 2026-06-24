package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/pquerna/otp/totp"
	_ "guineatrade.nhlstenden.com/src/1nit"
	"guineatrade.nhlstenden.com/src/database"
)

func TestExtractTOTP(t *testing.T) {
	w := httptest.NewRecorder()

	recoveryCode := "an invalid code lol"
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	result, err := ExtractTOTP(c)
	assert.Equal(t, result, "")
	assert.NotEqual(t, err, nil)
	assert.Equal(t, err.Error(), "can't find token in HTTP headers")

	c.Request.Header.Add("X-TOTP-Code", recoveryCode)

	result, err = ExtractTOTP(c)
	assert.Equal(t, result, recoveryCode)
	assert.Equal(t, err, nil)
}

func TestExtractRecoveryCode(t *testing.T) {
	w := httptest.NewRecorder()

	recoveryCode := "an invalid code lol"
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	result, err := ExtractRecoveryCode(c)
	assert.Equal(t, result, "")
	assert.NotEqual(t, err, nil)
	assert.Equal(t, err.Error(), "can't find token in HTTP headers")

	c.Request.Header.Add("X-Recovery-Code", recoveryCode)

	result, err = ExtractRecoveryCode(c)
	assert.Equal(t, result, recoveryCode)
	assert.Equal(t, err, nil)
}

func TestTOTPMiddleware(t *testing.T) {
	router := gin.Default()
	router.Use(TotpMiddlewareAuth())
	router.GET("/me", func(context *gin.Context) {
		context.Status(200)
	})

	user := database.CreateRandomUser()
	jwt, _ := GenerateToken(user)

	database.GetInstance().First(&user)

	type test struct {
		JWT            string
		TOTPCode       string
		HttpStatusCode int
		Result         string
	}

	tests := []test{
		{jwt, user.TotpSecret, http.StatusOK, ""},
		{jwt, "LSU5PZCBARGU63V5Y4CFQA5R5T", http.StatusUnauthorized, "{\"error\":\"MFA token is invalid\"}"},
		{jwt, "", http.StatusUnprocessableEntity, "{\"error\":\"No MFA supplied\"}"},
		{"", user.TotpSecret, http.StatusUnprocessableEntity, "{\"error\":\"Token expired\"}"},
	}
	for _, t2 := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/me", nil)
		if len(t2.JWT) > 0 {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t2.JWT))
		}
		code, _ := totp.GenerateCode(t2.TOTPCode, time.Now())
		if len(t2.TOTPCode) > 0 {
			req.Header.Set("X-TOTP-Code", code)
		}
		router.ServeHTTP(w, req)

		assert.Equal(t, t2.HttpStatusCode, w.Code)
		assert.Equal(t, t2.Result, w.Body.String())
	}
}
