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

func RegisterTOTP(c *gin.Context) {
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

func VerifyTOTP(c *gin.Context) {
	user, err := middleware.ExtractTokenUser(c)
	if err != nil {
		auth.SendError(c, http.StatusNotFound, err)
		return
	}
	passcode, err := middleware.ExtractTOTP(c)
	if err != nil {
		auth.SendError(c, http.StatusUnauthorized, err)
		return
	}

	database.GetInstance().First(&user)

	if !totp.Validate(passcode, user.TotpSecret) {
		auth.SendError(c, http.StatusUnauthorized, errors.New("invalid code"))
		return
	}

	c.Status(http.StatusOK)
}

func ResetTOTP(c *gin.Context) {
	user, err := middleware.ExtractTokenUser(c)
	if err != nil {
		auth.SendError(c, http.StatusNotFound, err)
		return
	}
	database.GetInstance().First(&user)

	passcode, err := middleware.ExtractTOTP(c)
	if err == nil {
		if !totp.Validate(passcode, user.TotpSecret) {
			auth.SendError(c, http.StatusUnauthorized, errors.New("invalid code"))
			return
		}
	} else {
		recoveryCode, err := middleware.ExtractRecoveryCode(c)
		if err != nil {
			c.Status(http.StatusUnauthorized)
		}
		if recoveryCode != user.RecoveryCode {
			auth.SendError(c, http.StatusUnauthorized, errors.New("recovery code does not match"))
			return
		}
	}
	user.TotpSecret = ""
	user.RecoveryCode = ""
	database.GetInstance().Save(&user)

	c.Status(http.StatusNoContent)
}
