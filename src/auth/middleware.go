package auth

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"guineatrade.nhlstenden.com/src/database"
)

var tokenLifeSpan int
var jwtSecret []byte

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s\n", err)
	}
	tokenLifeSpanString, lifespanExists := os.LookupEnv("JWT_TIMEOUT_HOURS")
	jwtSecretString, secretsExists := os.LookupEnv("JWT_SECRET_KEY")

	if !secretsExists || !lifespanExists {
		log.Fatal("JWT_TIMEOUT_HOURS and/or JWT_SECRET_KEY unset")
	}

	tokenLifeSpan, err = strconv.Atoi(tokenLifeSpanString)
	if err != nil {
		log.Fatal("JWT_TIMEOUT_HOURS is not a valid integer")
	}

	if len(jwtSecretString) < 64 {
		log.Fatal("JWT_SECRET_KEY is too short, minimum of 64 characters")
	}
	jwtSecret = []byte(jwtSecretString)
}

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := IsTokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

func GenerateToken(user *database.User) (string, error) {

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifeSpan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString([]byte(jwtSecret))
}

func IsTokenValid(c *gin.Context) error {
	tokenString, err := ExtractToken(c)
	if err != nil {
		return err
	}

	_, err = jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return err
	}

	return nil
}

func ExtractToken(c *gin.Context) (string, error) {
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) != 2 {
		return "", errors.New("can't find token in HTTP headers")
	}

	return strings.Split(bearerToken, " ")[1], nil
}

func ExtractTokenUser(c *gin.Context) (database.User, error) {
	tokenString, err := ExtractToken(c)
	if err != nil {
		return database.User{}, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return database.User{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return database.User{}, errors.New("token is invalid")
	}

	userId, ok := claims["user_id"].(float64)
	if !ok {
		return database.User{}, errors.New("supplied ID is not a valid integer")
	}

	var user database.User
	user.ID = uint(userId)
	if result := database.GetInstance().First(&user); result.Error != nil {
		return database.User{}, result.Error
	}

	return user, nil
}
