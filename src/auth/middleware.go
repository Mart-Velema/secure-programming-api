package auth

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
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
	claims["user_email"] = user.Email
	claims["user_name"] = user.Name
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespanInt)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString([]byte(jwtSecret))
}
