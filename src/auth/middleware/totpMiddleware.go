package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"guineatrade.nhlstenden.com/src/database"
)

func TotpMiddlewareAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := ExtractTokenUser(c)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}
		passcode, err := ExtractTOTP(c)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "No MFA supplied"})
			c.Abort()
			return
		}

		database.GetInstance().First(&user)

		if !totp.Validate(passcode, user.TotpSecret) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "MFA token is invalid"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func ExtractTOTP(c *gin.Context) (string, error) {
	TotpToken := c.Request.Header.Get("X-TOTP-Code")
	if len(TotpToken) == 0 {
		return "", errors.New("can't find token in HTTP headers")
	}

	return TotpToken, nil
}

func ExtractRecoveryCode(c *gin.Context) (string, error) {
	TotpToken := c.Request.Header.Get("X-Recovery-Code")
	if len(TotpToken) == 0 {
		return "", errors.New("can't find token in HTTP headers")
	}

	return TotpToken, nil
}
