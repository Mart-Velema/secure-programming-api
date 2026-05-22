package mfa

import (
	"crypto"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"guineatrade.nhlstenden.com/src/auth"
	"guineatrade.nhlstenden.com/src/auth/middleware"
	"guineatrade.nhlstenden.com/src/database"
)

type TotpCodes struct {
	Code         string `json:"code,omitempty"`
	RecoveryCode string `json:"recoveryCode,omitempty"`
}

func RegisterTOPT(c *gin.Context) {
	user, err := middleware.ExtractTokenUser(c)
	if err != nil {
		auth.SendError(c, http.StatusNotFound, err)
		return
	}

	randomBytes := rand.Text()
	recoveryCode := rand.Text() + rand.Text()

	hash := crypto.MD5.New()
	hash.Write([]byte(recoveryCode))
	recoveryHash := hex.EncodeToString(hash.Sum(nil))

	user.RecoveryCode = recoveryHash
	user.TotpSecret = randomBytes
	database.GetInstance().Save(&user)

	c.JSON(http.StatusOK, TotpCodes{
		Code:         randomBytes,
		RecoveryCode: recoveryCode,
	})
}

func VerifyTOPT(c *gin.Context) {

}

func ResetTOPT(c *gin.Context) {
	c.Status(http.StatusOK)
}
