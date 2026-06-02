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

const (
	backpackBaseUrl = "https://backpack.tf/api"
)

var (
	apiKey        string
	PricingCache  PricingDataCache
	CurrencyCache CurrencyDataCache
	itemCache     map[string]string
)

var client = http.Client{
	Transport: &http.Transport{
		MaxIdleConns:    30,
		MaxConnsPerHost: 10,
		IdleConnTimeout: 30 * time.Second,
		TLSClientConfig: &tls.Config{
			ServerName:         "backpack.tf",
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
	if conn == nil {
		log.Fatal(err)
	}
	defer func(conn *tls.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

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
		installItemCache()
		for {
			updatePriceCache()
			updateCurrencyCache()
			now := time.Now().Truncate(time.Hour)

			timeTillNextUpdate := 6 - (now.Hour() % 6)
			now = now.Add(time.Duration(timeTillNextUpdate) * time.Hour)

			nextRefresh := time.Until(now)
			log.Printf("Next pricing update: %s", now.String())

			time.Sleep(nextRefresh)
		}
	}()
}

func getPrice() (*pricingData, error) {
	var pricingResponse pricingData
	response, err := client.Get(fmt.Sprintf("%s/IGetPrices/v4?key=%s", backpackBaseUrl, apiKey))
	if err != nil {
		return &pricingResponse, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(response.Body)

	if response.StatusCode != 200 {
		return &pricingResponse, fmt.Errorf("unable to get current pricing: %d", response.StatusCode)
	}

	decoder := json.NewDecoder(response.Body)
	for {
		err := decoder.Decode(&pricingResponse)

		if err == io.EOF {
			break
		} else if err != nil {
			return &pricingResponse, err
		}
	}

	return &pricingResponse, nil
}

func getCurrency() (*currencyData, error) {
	var currencyResponse currencyData
	response, err := client.Get(fmt.Sprintf("%s/IGetCurrencies/v1?key=%s", backpackBaseUrl, apiKey))
	if err != nil {
		return &currencyResponse, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(response.Body)

	if response.StatusCode != 200 {
		return &currencyResponse, fmt.Errorf("unable to get current currency conversions: %d", response.StatusCode)
	}

	decoder := json.NewDecoder(response.Body)
	for {
		err := decoder.Decode(&currencyResponse)

		if err == io.EOF {
			break
		} else if err != nil {
			return &currencyResponse, err
		}
	}

	return &currencyResponse, nil
}

func updatePriceCache() {
	priceResult, err := getPrice()
	if err != nil {
		log.Println(err)
		return
	}
	priceCache, err := priceResult.toCache()
	if err != nil {
		log.Println(err)
		log.Println("Using old cache")
		return
	}
	PricingCache = *priceCache
	log.Printf("Updated Price cache on %s", time.Now().String())
}

func updateCurrencyCache() {
	priceResult, err := getCurrency()
	if err != nil {
		log.Println(err)
		return
	}
	currencyCache := priceResult.toCache()

	CurrencyCache = *currencyCache
	log.Printf("Updated Currency cache on %s", time.Now().String())
}

func installItemCache() {
	content, err := os.ReadFile("./item-icons.json")
	if err != nil {
		log.Fatal(err)
	}

	var items = make(map[string]string)
	err = json.Unmarshal(content, &items)
	if err != nil {
		log.Fatal(err)
	}

	itemCache = items
}
