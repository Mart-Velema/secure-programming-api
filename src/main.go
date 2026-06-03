package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"guineatrade.nhlstenden.com/src/auth"
	"guineatrade.nhlstenden.com/src/auth/mfa"
	"guineatrade.nhlstenden.com/src/auth/middleware"
	"guineatrade.nhlstenden.com/src/backpack"
	"guineatrade.nhlstenden.com/src/database"
	"guineatrade.nhlstenden.com/src/steam"
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

	steamPublic := router.Group("/api/v1/steam")
	{
		steamPublic.GET("/status", steam.GetBotStatus)
	}

	apiRestricted := router.Group("/api/v1")
	apiRestricted.Use(middleware.JwtAuthMiddleware())
	{
		authGroup := apiRestricted.Group("/auth")
		{
			authGroup.POST("/logout", auth.Logout)
			authGroup.POST("/logout/all", auth.LogoutAll)
			authGroup.GET("/me", auth.Me)
			authGroup.PATCH("/me", auth.UpdatePassword)

			multifactorAuthGroup := authGroup.Group("/mfa")
			{
				multifactorAuthGroup.POST("/totp/register", mfa.RegisterTOTP)
				multifactorAuthGroup.POST("/totp/verify", mfa.VerifyTOTP)
				multifactorAuthGroup.DELETE("/totp/reset", mfa.ResetTOTP)
			}
		}
		backpackGroup := apiRestricted.Group("/backpack")
		{
			backpackGroup.GET("/prices", backpack.GetPrices)
			backpackGroup.GET("/prices/:item", backpack.GetItemDetails)
			backpackGroup.GET("/currency", backpack.GetCurrencies)
		}
	}

	err := router.Run(fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatal(err)
	}
}
