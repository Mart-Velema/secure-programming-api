package auth

import (
	"net/http"
	"regexp"
	"time"

	"github.com/cgholdings/go-common/database/encryption"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"guineatrade.nhlstenden.com/src/auth/middleware"
	"guineatrade.nhlstenden.com/src/database"
)

type registerUser struct {
	Email          string `json:"email"`
	Name           string `json:"name"`
	Password       string `json:"password"`
	PasswordVerify string `json:"passwordVerify"`
}

type patchUser struct {
	Email             string `json:"email"`
	CurrentPassword   string `json:"currentPassword"`
	NewPassword       string `json:"newPassword"`
	NewPasswordVerify string `json:"newPasswordVerify"`
}

type patchSteam struct {
	TradeUrl string `json:"tradeUrl"`
	SteamId  uint64 `json:"steamId"`
}

func (user *registerUser) toDatabaseRecord() database.User {
	return database.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
}

type Tokens struct {
	JWT     string `json:"jwt,omitempty"`
	Refresh string `json:"refresh,omitempty"`
}

func Register(c *gin.Context) {
	var postRegister registerUser
	if err := c.ShouldBindJSON(&postRegister); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid input data"})
		return
	}
	if !isEmailValid(postRegister.Email) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Email is not valid"})
		return
	}
	if postRegister.Password != postRegister.PasswordVerify {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Passwords do not match"})
		return
	}

	if result := database.GetInstance().Create(new(postRegister.toDatabaseRecord())); result.Error != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}
	c.Status(http.StatusCreated)
}

func Login(c *gin.Context) {
	var user database.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid input data"})
		return
	}

	if result := database.GetInstance().
		Where("email_hash = ?", encryption.Hash(user.Email)).
		Where("password = ?", encryption.Hash(user.Password)).
		First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User does not exist"})
		return
	}

	JWT, err := middleware.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to generate token"})
		return
	}
	refreshToken, err := middleware.GenerateRefreshToken(&user, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, Tokens{
		JWT:     JWT,
		Refresh: refreshToken,
	})
}

func Refresh(c *gin.Context) {
	var refreshToken Tokens
	if err := c.ShouldBindJSON(&refreshToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid input data"})
		return
	}

	var token database.RefreshToken
	if result := database.GetInstance().
		Joins("User").
		Where("refresh_tokens.token_hash = ?", encryption.Hash(refreshToken.Refresh)).
		First(&token); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token does not exist"})
		return
	}

	if time.Now().After(token.ExpiresOn) {
		database.GetInstance().Delete(&token)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is expired"})
		return
	}
	if token.Nonce != middleware.GenerateTokenNonce(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is used from another location"})
		return
	}

	jwt, err := middleware.GenerateToken(&token.User)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to generate token"})
		return
	}

	c.JSON(http.StatusOK, &Tokens{
		JWT: jwt,
	})
}

func Logout(c *gin.Context) {
	var refreshToken Tokens
	if err := c.ShouldBindJSON(&refreshToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid input data"})
		return
	}

	if result := database.GetInstance().Delete(&refreshToken.Refresh); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cannot find logout credentials"})
		return
	}

	c.Status(http.StatusNoContent)
}

func LogoutAll(c *gin.Context) {
	user, err := middleware.ExtractTokenUser(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Token expired"})
		return
	}

	if result := database.GetInstance().
		Where("user_id = ?", &user.ID).
		Delete(&database.RefreshToken{}); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cannot find logout credentials"})
		return
	}

	c.Status(http.StatusNoContent)
}

func Me(c *gin.Context) {
	user, err := middleware.ExtractTokenUser(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Token expired"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdatePassword(c *gin.Context) {
	user, err := middleware.ExtractTokenUser(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Token expired"})
		return
	}

	var requestUser patchUser
	if err := c.ShouldBindJSON(&requestUser); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid input data"})
		return
	}

	if requestUser.NewPassword != requestUser.NewPasswordVerify {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	if user.HasMFAEnabled() {
		totpToken, err := middleware.ExtractTOTP(c)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No MFA supplied"})
			return
		}
		if totp.Validate(totpToken, user.TotpSecret) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid MFA"})
			return
		}
	}

	user.Password = requestUser.NewPassword
	database.GetInstance().Select("password").Save(&user)

	if result := database.GetInstance().
		Where("user_id = ?", user.ID).
		Delete(&database.RefreshToken{}); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cannot find logout credentials"})
		return
	}

	c.Status(http.StatusAccepted)
}

func UpdateSteam(c *gin.Context) {
	user, err := middleware.ExtractTokenUser(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Token expired"})
		return
	}

	var steamPatch patchSteam
	if err = c.ShouldBindJSON(&steamPatch); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid input data"})
		return
	}

	user.TradeUrl = steamPatch.TradeUrl
	user.SteamId = steamPatch.SteamId

	if result := database.
		GetInstance().
		Select("steam_id", "trade_url").
		Save(&user); result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to update Steam credentials"})
		return
	}

	c.Status(http.StatusNoContent)
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}
