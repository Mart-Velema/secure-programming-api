package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"guineatrade.nhlstenden.com/src/database"
)

func totpMiddlewareAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := ExtractTokenUser(c)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}
		passcode, err := ExtractTOTP(c)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		database.GetInstance().First(&user)

		if totp.Validate(passcode, user.TotpSecret) {
			c.Status(http.StatusUnauthorized)
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
