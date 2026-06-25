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
	"guineatrade.nhlstenden.com/src/inventory"
	"guineatrade.nhlstenden.com/src/steam"
	"guineatrade.nhlstenden.com/src/stripe"
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
	_ = database.GetInstance()
}

func main() {
	fmt.Printf("Running locally on localhost:%s\n", os.Getenv("PORT"))

	router := gin.Default()

	apiPublic := router.Group("/api/v1")
	{
		authGroup := apiPublic.Group("/auth")
		{
			authGroup.POST("/register", auth.Register)
			authGroup.POST("/login", auth.Login)
			authGroup.POST("/refresh", auth.Refresh)
		}

		apiPublic.POST("/stripe/webhook", stripe.Webhook)
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
			authGroup.PATCH("/steam", auth.UpdateSteam)

			multifactorAuthGroup := authGroup.Group("/mfa")
			{
				multifactorAuthGroup.POST("/totp/register", mfa.RegisterTOTP)
				multifactorAuthGroup.POST("/totp/verify", mfa.VerifyTOTP)
				multifactorAuthGroup.DELETE("/totp/reset", mfa.ResetTOTP)
			}
		}
		apiRestricted.GET("/backpack/prices", backpack.GetPrices)

		steamGroup := apiRestricted.Group("/steam")
		{
			steamGroup.GET("/inventory", steam.GetBotInventory)
			steamGroup.GET("/stock", inventory.GetSteamBotStock)
		}

		userGroup := apiRestricted.Group("/user")
		{
			userGroup.GET("/inventory", inventory.GetInventory)
			userGroup.GET("/stock", inventory.GetUserStock)
			userGroup.GET("/trade/status", steam.GetTradeStatus)
		}
		stripeGroup := apiRestricted.Group("/stripe")
		{
			stripeGroup.POST("/create", stripe.CreatePaymentSession)
		}
	}

	err := router.Run(fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatal(err)
	}
}
