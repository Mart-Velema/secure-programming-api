package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"guineatrade.nhlstenden.com/src/database"
)

func HelloWorld(c *gin.Context) {
	c.String(200, "%s", "Hello, world!")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s\n", err)
	}
	fmt.Printf("Running locally on localhost:%s\n", os.Getenv("PORT"))

	_ = database.GetInstance()

	router := gin.Default()

	router.GET("/api/v1", HelloWorld)

	err = router.Run(fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatal(err)
	}
}
