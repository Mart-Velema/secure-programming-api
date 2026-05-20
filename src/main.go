package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"guineatrade.nhlstenden.com/src/auth"
	"guineatrade.nhlstenden.com/src/auth/mfa"
	"guineatrade.nhlstenden.com/src/database"
)

func HelloWorld(c *gin.Context) {
	c.String(200, "%s", "Hello, world!")
}

var Version string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s\n", err)
	}
	if len(os.Args) >= 2 && os.Args[1] == "--seed" {
		err := os.Remove(os.Getenv("SQLITE_FILE_LOCATION"))
		if err != nil {
			log.Println(err)
		}
		database.Seed()
	}
	_ = database.GetInstance()
}

func main() {
	fmt.Printf("Running locally on localhost:%s\n", os.Getenv("PORT"))

	router := gin.Default()

	apiPublic := router.Group("/api/v1/auth")
	{
		apiPublic.POST("/register", auth.Register)
		apiPublic.POST("/login", auth.Login)
		apiPublic.POST("/refresh", auth.Refresh)
	}

	apiRestricted := router.Group("/api/v1")
	apiRestricted.Use(auth.JwtAuthMiddleware())
	{
		authGroup := apiRestricted.Group("/auth")
		{
			authGroup.POST("/logout", auth.Logout)
			authGroup.POST("/logout/all", auth.LogoutAll)
			authGroup.GET("/me", auth.Me)

			multifactorAuthGroup := authGroup.Group("/mfa")
			{
				multifactorAuthGroup.POST("/sms/send", mfa.SendSMS)
				multifactorAuthGroup.POST("/sms/verify", mfa.VerifySMS)
			}
		}
	}

	err := router.Run(fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatal(err)
	}
}
