package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	_ "guineatrade.nhlstenden.com/src/1nit"
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
