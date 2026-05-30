package backpack

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var apiKey string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s\n", err)
	}
	envApiKey, apiKeyExists := os.LookupEnv("BACKPACK_API_KEY")

	if !apiKeyExists {
		log.Fatal("BACKPACK_API_KEY is unset")
	}

	apiKey = envApiKey

	go func() {
		for {
			getPrice()
			getCurrency()
			time.Sleep(time.Hour * 24)
		}
	}()
}

func getPrice() {
	//	TODO: Cache these
	//  TODO: Use proper remote URL instead of local testing URL
	response, err := http.Get(fmt.Sprintf("http://localhost:8080/api/IGetPrices/v4?key=%s", apiKey))
	if err != nil {
		log.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(response.Body)

	if response.StatusCode != 200 {
		log.Printf("Unable to get current pricing: %d", response.StatusCode)
		return
	}

	decoder := json.NewDecoder(response.Body)
	var pricingResponse PricingData
	for {
		err := decoder.Decode(&pricingResponse)

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error decoding:", err)
			return
		}
	}

	fmt.Printf("Received item: %+v\n\n", pricingResponse)

}

func getCurrency() {
	//	TODO: Cache these
	//  TODO: Use proper remote URL instead of local testing URL
	response, err := http.Get(fmt.Sprintf("http://localhost:8080/api/IGetCurrencies/v1?key=%s", apiKey))
	if err != nil {
		log.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(response.Body)

	if response.StatusCode != 200 {
		log.Printf("Unable to get current currency conversion: %d", response.StatusCode)
		return
	}

	decoder := json.NewDecoder(response.Body)
	var currencyResponse CurrencyData
	for {
		err := decoder.Decode(&currencyResponse)

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error decoding:", err)
			return
		}
	}

	fmt.Printf("Received item: %+v\n", currencyResponse)
}
