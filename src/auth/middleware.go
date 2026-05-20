package auth

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"guineatrade.nhlstenden.com/src/database"
)

func GenerateToken(user *database.User) (string, error) {
	tokenLifespan, lifespanExists := os.LookupEnv("JWT_TIMEOUT_HOURS")
	jwtSecret, secretExists := os.LookupEnv("JWT_SECRET_KEY")

	if !secretExists || !lifespanExists {
		log.Fatal("JWT_TIMEOUT_HOURS and/or JWT_SECRET_KEY unset")
	}
	tokenLifespanInt, err := strconv.Atoi(tokenLifespan)
	if err != nil {
		log.Fatal("JWT_TIMEOUT_HOURS is not a valid integer")
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespanInt)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString([]byte(jwtSecret))
}

func ExtractToken(c *gin.Context) (string, error) {
	token := c.Query("token")
	if token != "" {
		return token, nil
	}

	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1], nil
	}

	return "", errors.New("can't find token in HTTP headers")
}

func ExtractTokenUser(c *gin.Context) (database.User, error) {
	tokenString, err := ExtractToken(c)
	if err != nil {
		return database.User{}, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return database.User{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return database.User{}, errors.New("token is invalid")
	}

	userId, err := strconv.Atoi(claims["user_id"].(string))
	if err != nil {
		return database.User{}, err
	}

	var user database.User
	user.ID = uint(userId)
	if result := database.GetInstance().First(&user); result.Error != nil {
		return database.User{}, result.Error
	}

	return user, nil
}
