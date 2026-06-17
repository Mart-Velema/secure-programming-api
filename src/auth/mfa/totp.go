package mfa

import (
	"crypto"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
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
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Token expired"})
		return
	}

	randomBytes := rand.Text()
	recoveryCode := rand.Text() + rand.Text()

	hash := crypto.MD5.New()
	hash.Write([]byte(recoveryCode))
	recoveryHash := hex.EncodeToString(hash.Sum(nil))

	user.RecoveryCode = recoveryHash
	user.TotpSecret = randomBytes
	database.GetInstance().Select("totp_secret", "recovery_code").Save(&user)

	c.JSON(http.StatusOK, TotpCodes{
		Code:         randomBytes,
		RecoveryCode: recoveryHash,
	})
}

func VerifyTOTP(c *gin.Context) {
	user, err := middleware.ExtractTokenUser(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token expired"})
		return
	}
	passcode, err := middleware.ExtractTOTP(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No MFA supplied"})
		return
	}

	database.GetInstance().First(&user)

	if !totp.Validate(passcode, user.TotpSecret) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid MFA"})
		return
	}

	c.Status(http.StatusOK)
}

func ResetTOTP(c *gin.Context) {
	user, err := middleware.ExtractTokenUser(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Token expired"})
		return
	}
	database.GetInstance().First(&user)

	passcode, err := middleware.ExtractTOTP(c)
	if err == nil {
		if !totp.Validate(passcode, user.TotpSecret) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid MFA"})
			return
		}
	} else {
		recoveryCode, err := middleware.ExtractRecoveryCode(c)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "No recovery or TOTP code supplied"})
		}
		if recoveryCode != user.RecoveryCode {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid recovery code"})
			return
		}
	}

	user.TotpSecret = ""
	user.RecoveryCode = ""
	database.GetInstance().Select("totp_secret", "recovery_code").Save(&user)

	c.Status(http.StatusNoContent)
}
