package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"guineatrade.nhlstenden.com/src/database"
)

type Passcode struct {
	Code string `json:"code"`
}

func totpMiddlewareAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := ExtractTokenUser(c)
		if err != nil {
			c.Status(http.StatusNotFound)
			c.Abort()
			return
		}
		var passcode Passcode
		err = c.ShouldBindJSON(&passcode)
		if err != nil {
			c.Status(http.StatusNotFound)
			c.Abort()
			return
		}

		database.GetInstance().First(&user)

		if totp.Validate(passcode.Code, user.TotpSecret) {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
