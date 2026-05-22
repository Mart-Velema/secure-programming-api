package mfa

import (
	"crypto"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"guineatrade.nhlstenden.com/src/auth"
	"guineatrade.nhlstenden.com/src/auth/middleware"
	"guineatrade.nhlstenden.com/src/database"
)

type TotpCodes struct {
	Code         string `json:"code,omitempty"`
	RecoveryCode string `json:"recovery,omitempty"`
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
		RecoveryCode: recoveryHash,
	})
}

func VerifyTOPT(c *gin.Context) {
	user, err := middleware.ExtractTokenUser(c)
	if err != nil {
		auth.SendError(c, http.StatusNotFound, err)
		return
	}
	var passcode TotpCodes
	err = c.ShouldBindJSON(&passcode)
	if err != nil {
		auth.SendError(c, http.StatusBadRequest, err)
		return
	}

	database.GetInstance().First(&user)

	if !totp.Validate(passcode.Code, user.TotpSecret) {
		auth.SendError(c, http.StatusUnauthorized, errors.New("invalid code"))
		return
	}

	c.Status(http.StatusOK)
}

func ResetTOPT(c *gin.Context) {
	user, err := middleware.ExtractTokenUser(c)
	if err != nil {
		auth.SendError(c, http.StatusNotFound, err)
		return
	}
	var passcode TotpCodes
	err = c.ShouldBindJSON(&passcode)
	if err != nil {
		auth.SendError(c, http.StatusBadRequest, err)
		return
	}

	database.GetInstance().First(&user)
	if passcode.RecoveryCode != user.RecoveryCode {
		auth.SendError(c, http.StatusUnauthorized, errors.New("recovery code does not match"))
		return
	}

	user.TotpSecret = ""
	user.RecoveryCode = ""
	database.GetInstance().Save(&user)

	c.Status(http.StatusNoContent)
}
