package backpack

import (
	"crypto/tls"
	"crypto/x509"
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

var client = http.Client{
	Transport: &http.Transport{
		MaxIdleConns:    30,
		MaxConnsPerHost: 10,
		IdleConnTimeout: 30 * time.Second,
		TLSClientConfig: &tls.Config{
			ServerName:         "backpac.tf",
			InsecureSkipVerify: false,
		},
	},
}

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

	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := tls.Dial("tcp", "backpack.tf:443", &tls.Config{
		RootCAs:    rootCAs,
		ServerName: "backpack.tf",
	})
	defer conn.Close()

	connState := conn.ConnectionState()
	chain := connState.PeerCertificates

	intermediates := x509.NewCertPool()
	for _, cert := range chain {
		intermediates.AddCert(cert)
	}

	opts := x509.VerifyOptions{
		Roots:         rootCAs,
		DNSName:       "backpack.tf",
		Intermediates: intermediates,
	}
	_, err = chain[0].Verify(opts)
	if err != nil {
		log.Fatal(err)
	}

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
	response, err := client.Get(fmt.Sprintf("http://localhost:8080/api/IGetPrices/v4?key=%s", apiKey))
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
			log.Println("Error decoding:", err)
			return
		}
	}

	fmt.Printf("Received item: %+v\n\n", pricingResponse)

}

func getCurrency() {
	//	TODO: Cache these
	//  TODO: Use proper remote URL instead of local testing URL
	response, err := client.Get(fmt.Sprintf("http://localhost:8080/api/IGetCurrencies/v1?key=%s", apiKey))
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
			log.Println("Error decoding:", err)
			return
		}
	}

	fmt.Printf("Received item: %+v\n", currencyResponse)
}
