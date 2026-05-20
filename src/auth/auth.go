package auth

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/cgholdings/go-common/database/encryption"
	"github.com/gin-gonic/gin"
	"guineatrade.nhlstenden.com/src/database"
)

type registerUser struct {
	Email          string `json:"email"`
	Name           string `json:"name"`
	Password       string `json:"password"`
	PasswordVerify string `json:"passwordVerify"`
	PhoneNumber    string `json:"tel"`
}

func (user *registerUser) toDatabaseRecord() database.User {
	return database.User{
		Name:        user.Name,
		Email:       user.Email,
		Password:    user.Password,
		PhoneNumber: user.PhoneNumber,
	}
}

type loginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Tokens struct {
	JWT     string `json:"jwt,omitempty"`
	Refresh string `json:"refresh,omitempty"`
}

func Register(c *gin.Context) {
	var postRegister registerUser
	if err := c.BindJSON(&postRegister); err != nil {
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
	var loggedinUser loginUser
	if err := c.BindJSON(&loggedinUser); err != nil {
		SendError(c, http.StatusUnprocessableEntity, err)
		return
	}

	var user database.User
	if result := database.GetInstance().
		Where("email_hash = ?", encryption.Hash(loggedinUser.Email)).
		Where("password = ?", encryption.Hash(loggedinUser.Password)).
		First(&user); result.Error != nil {
		SendError(c, http.StatusNotFound, result.Error)
		return
	}

	JWT, err := GenerateToken(&user)
	if err != nil {
		SendError(c, http.StatusUnauthorized, err)
		return
	}
	refreshToken, err := GenerateRefreshToken(&user)
	if err != nil {
		SendError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, Tokens{
		JWT:     JWT,
		Refresh: refreshToken,
	})
}

func Logout(c *gin.Context) {
	c.Status(http.StatusOK)
}

func Me(c *gin.Context) {
	user, err := ExtractTokenUser(c)
	if err != nil {
		SendError(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, user)
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
