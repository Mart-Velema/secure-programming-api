package auth

import (
	"errors"
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
		SendError(c, http.StatusUnprocessableEntity, err)
		return
	}
	if !isEmailValid(postRegister.Email) {
		SendError(c, http.StatusUnprocessableEntity, errors.New("supplied email address is not valid"))
		return
	}
	if postRegister.Password != postRegister.PasswordVerify {
		SendError(c, http.StatusUnprocessableEntity, errors.New("supplied passwords do not match"))
		return
	}

	if result := database.GetInstance().Create(new(postRegister.toDatabaseRecord())); result.Error != nil {
		SendError(c, http.StatusInternalServerError, result.Error)
		return
	}
	c.Status(http.StatusCreated)
}

func Login(c *gin.Context) {
	var user database.User
	if err := c.ShouldBindJSON(&user); err != nil {
		SendError(c, http.StatusUnprocessableEntity, err)
		return
	}

	if result := database.GetInstance().
		Where("email_hash = ?", encryption.Hash(user.Email)).
		Where("password = ?", encryption.Hash(user.Password)).
		First(&user); result.Error != nil {
		SendError(c, http.StatusNotFound, result.Error)
		return
	}

	JWT, err := middleware.GenerateToken(&user)
	if err != nil {
		SendError(c, http.StatusUnauthorized, err)
		return
	}
	refreshToken, err := middleware.GenerateRefreshToken(&user, c)
	if err != nil {
		SendError(c, http.StatusInternalServerError, err)
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
		SendError(c, http.StatusUnprocessableEntity, err)
		return
	}

	var token database.RefreshToken
	if result := database.GetInstance().
		Joins("User").
		Where("refresh_tokens.token_hash = ?", encryption.Hash(refreshToken.Refresh)).
		First(&token); result.Error != nil {
		SendError(c, http.StatusNotFound, result.Error)
		return
	}

	if time.Now().After(token.ExpiresOn) {
		database.GetInstance().Delete(&token)
		SendError(c, http.StatusUnauthorized, errors.New("login expired"))
		return
	}
	if token.Nonce != middleware.GenerateTokenNonce(c) {
		SendError(c, http.StatusUnauthorized, errors.New("logged in from another location"))
		return
	}

	jwt, err := middleware.GenerateToken(&token.User)
	if err != nil {
		SendError(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, &Tokens{
		JWT: jwt,
	})
}

func Logout(c *gin.Context) {
	var refreshToken Tokens
	if err := c.ShouldBindJSON(&refreshToken); err != nil {
		SendError(c, http.StatusUnprocessableEntity, err)
		return
	}

	if result := database.GetInstance().Delete(&refreshToken.Refresh); result.Error != nil {
		SendError(c, http.StatusNotFound, result.Error)
		return
	}

	c.Status(http.StatusNoContent)
}

func LogoutAll(c *gin.Context) {
	user, err := middleware.ExtractTokenUser(c)
	if err != nil {
		SendError(c, http.StatusNotFound, err)
		return
	}

	if result := database.GetInstance().
		Where("user_id = ?", &user.ID).
		Delete(&database.RefreshToken{}); result.Error != nil {
		SendError(c, http.StatusNotFound, result.Error)
		return
	}

	c.Status(http.StatusNoContent)
}

func Me(c *gin.Context) {
	user, err := middleware.ExtractTokenUser(c)
	if err != nil {
		SendError(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdatePassword(c *gin.Context) {
	var requestUser patchUser
	if err := c.ShouldBindJSON(&requestUser); err != nil {
		SendError(c, http.StatusUnprocessableEntity, err)
		return
	}

	var user database.User
	if result := database.GetInstance().
		Where("email_hash = ?", encryption.Hash(requestUser.Email)).
		Where("password = ?", encryption.Hash(requestUser.CurrentPassword)).
		First(&user); result.Error != nil {
		SendError(c, http.StatusNotFound, result.Error)
		return
	}

	if requestUser.NewPassword != requestUser.NewPasswordVerify {
		SendError(c, http.StatusBadRequest, errors.New("passwords do not match"))
		return
	}

	if user.HasMFAEnabled() {
		totpToken, err := middleware.ExtractTOTP(c)
		if err != nil {
			SendError(c, http.StatusNotFound, err)
			return
		}
		if totp.Validate(totpToken, user.TotpSecret) {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}
	}

	user.Password = requestUser.NewPassword
	database.GetInstance().Select("password").Save(&user)

	if result := database.GetInstance().
		Where("user_id = ?", user.ID).
		Delete(&database.RefreshToken{}); result.Error != nil {
		SendError(c, http.StatusNotFound, result.Error)
		return
	}

	c.Status(http.StatusAccepted)
}

func UpdateSteam(c *gin.Context) {
	user, err := middleware.ExtractTokenUser(c)
	if err != nil {
		SendError(c, http.StatusNotFound, err)
		return
	}

	var steamPatch patchSteam
	if err := c.ShouldBindJSON(&steamPatch); err != nil {
		SendError(c, http.StatusUnprocessableEntity, err)
		return
	}

	user.TradeUrl = steamPatch.TradeUrl
	user.SteamId = steamPatch.SteamId

	if result := database.
		GetInstance().
		Select("steam_id", "trade_url").
		Save(&user); result.Error != nil {
		SendError(c, http.StatusUnprocessableEntity, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func SendError(c *gin.Context, statusCode int, err error) {
	type genericError struct {
		Message any `json:"message"`
	}

	c.JSON(statusCode, genericError{
		Message: err.Error(),
	})
}
