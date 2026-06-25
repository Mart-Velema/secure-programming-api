package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cgholdings/go-common/database/encryption"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang-jwt/jwt/v5"
	_ "guineatrade.nhlstenden.com/src/1nit"
	"guineatrade.nhlstenden.com/src/database"
)

func TestGenerateToken(t *testing.T) {
	type jwtHeader struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}

	type jwtPayload struct {
		Authorized bool  `json:"authorized"`
		UserId     uint  `json:"user_id"`
		Exp        int32 `json:"exp"`
	}

	user := database.CreateRandomUser()
	userJwt, err := GenerateToken(user)

	assert.Equal(t, err, nil)

	splits := strings.Split(userJwt, ".")
	assert.Equal(t, len(splits), 3)

	header, _ := base64.StdEncoding.DecodeString(splits[0])
	payload, _ := base64.StdEncoding.DecodeString(splits[1])

	var jsonHeader jwtHeader
	err = json.Unmarshal(header, &jsonHeader)
	assert.Equal(t, err, nil)
	assert.Equal(t, jsonHeader.Alg, "HS512")
	assert.Equal(t, jsonHeader.Typ, "JWT")

	var jsonPayload jwtPayload
	err = json.Unmarshal(payload, &jsonPayload)
	assert.Equal(t, err, nil)

	expTimestamp := time.Unix(int64(jsonPayload.Exp), 0)
	jwtTimeoutMinuteString := os.Getenv("JWT_TIMEOUT_MINUTES")
	jwtTimeoutMinutes, _ := strconv.Atoi(jwtTimeoutMinuteString)

	assert.Equal(t, jsonPayload.Authorized, true)
	assert.Equal(t, jsonPayload.UserId, user.ID)
	assert.Equal(t, time.Now().Before(expTimestamp), true)
	assert.Equal(t, time.Now().
		Truncate(time.Minute).
		Add(time.Minute*time.Duration(jwtTimeoutMinutes)), expTimestamp.Truncate(time.Minute))

}

func TestExtractToken(t *testing.T) {
	user := database.CreateRandomUser()
	userJwt, _ := GenerateToken(user)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	result, err := ExtractToken(c)

	assert.Equal(t, result, "")
	assert.NotEqual(t, err, nil)
	assert.Equal(t, err.Error(), "can't find token in HTTP headers")

	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userJwt))

	result, err = ExtractToken(c)
	assert.Equal(t, result, userJwt)
	assert.Equal(t, err, nil)

	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer%s", userJwt))

	result, err = ExtractToken(c)
	assert.Equal(t, result, "")
	assert.NotEqual(t, err, nil)
	assert.Equal(t, err.Error(), "can't find token in HTTP headers")
}

func TestJwtAuthMiddleware(t *testing.T) {
	type test struct {
		JWT        string
		StatusCode int
		Result     string
	}
	router := gin.Default()
	router.Use(JwtAuthMiddleware())
	router.GET("/jwt", func(context *gin.Context) {
		context.Status(200)
	})

	user := database.CreateRandomUser()
	userJwt, _ := GenerateToken(user)

	tests := []test{
		{userJwt, http.StatusOK, ""},
		{"", http.StatusUnauthorized, "{\"error\":\"Invalid token\"}"},
	}
	for _, t2 := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/jwt", nil)
		if len(t2.JWT) > 0 {
			req.Header.Set("Authorization", fmt.Sprintf("bearer %s", userJwt))
		}
		router.ServeHTTP(w, req)

		assert.Equal(t, t2.StatusCode, w.Code)
		assert.Equal(t, t2.Result, w.Body.String())
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	user := database.CreateRandomUser()
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("User-Agent", "unit-test/tokenNonce-1")

	refreshToken, _ := GenerateRefreshToken(user, c)

	var dbToken = database.RefreshToken{
		Token: refreshToken,
	}
	database.GetInstance().Find(&dbToken)

	refreshLifeSpanString := os.Getenv("JWT_REFRESH_DAYS")
	refreshLifeSpan, _ = strconv.Atoi(refreshLifeSpanString)

	timeoutTime := time.Now().Truncate(time.Hour).Add(time.Hour * time.Duration(24*refreshLifeSpan)).Unix()

	assert.Equal(t, user.ID, dbToken.UserID)
	assert.Equal(t, encryption.Hash(c.ClientIP()+"unit-test/tokenNonce-1"), dbToken.Nonce)
	assert.Equal(t, timeoutTime, dbToken.ExpiresOn.Truncate(time.Hour).Unix())
}

func TestGenerateTokenNonce(t *testing.T) {
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("User-Agent", "unit-test/tokenNonce-1")

	nonce := GenerateTokenNonce(c)

	assert.Equal(t, encryption.Hash(c.ClientIP()+"unit-test/tokenNonce-1"), nonce)
	assert.NotEqual(t, encryption.Hash(c.ClientIP()+"unit-test/tokenNonce-2"), nonce)
}

func TestExtractTokenUser(t *testing.T) {
	user := database.CreateRandomUser()

	type test struct {
		JWT    jwt.MapClaims
		Result string
	}

	tests := []test{
		{
			jwt.MapClaims{
				"authorized": true,
				"exp":        time.Now().Add(time.Hour * -1).Unix(),
				"user_id":    1,
			},
			"token is expired",
		},
		{
			jwt.MapClaims{
				"authorized": true,
				"exp":        time.Now().Add(time.Hour * 1).Unix(),
				"user_id":    "notanumber",
			},
			"supplied ID is not a valid integer",
		},
		{
			jwt.MapClaims{
				"authorized": true,
				"exp":        time.Now().Add(time.Hour * 1).Unix(),
				"user_id":    999,
			},
			"record not found",
		},
		{
			jwt.MapClaims{
				"authorized": true,
				"exp":        time.Now().Add(time.Hour * 1).Unix(),
				"user_id":    user.ID,
			},
			"",
		},
	}

	for _, t2 := range tests {
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", generateTestToken(t2.JWT)))

		result, err := ExtractTokenUser(c)

		if t2.Result == "" {
			assert.Equal(t, err, nil)
			assert.Equal(t, result, user)
		} else {
			assert.Equal(t, result, nil)
			assert.NotEqual(t, err, nil)
			assert.Equal(t, strings.Contains(err.Error(), t2.Result), true)
		}
	}
}

func generateTestToken(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		panic(err)
	}
	return tokenString
}
