package backpack

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	backpackBaseUrl = "https://backpack.tf/api"
)

var (
	apiKey        string
	PricingCache  PricingDataCache
	itemCache     map[string]itemConstants
	unusualCache  map[string]string
	defindexCache map[string]uint32
)

type itemConstants struct {
	Url            string `json:"url"`
	MarketHashName string `json:"marketHashName"`
}

var client *http.Client

func init() {
	envApiKey, apiKeyExists := os.LookupEnv("BACKPACK_API_KEY")
	envApiHash, apiHashExists := os.LookupEnv("BACKPACK_API_HASH")

	if !apiKeyExists || !apiHashExists {
		log.Fatal("BACKPACK_API_KEY or BACKPACK_API_HASH is unset")
	}
	client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    30,
			MaxConnsPerHost: 10,
			IdleConnTimeout: 30 * time.Second,
			TLSClientConfig: &tls.Config{
				ServerName:         "backpack.tf",
				InsecureSkipVerify: false,
				VerifyConnection: func(cs tls.ConnectionState) error {
					if len(cs.PeerCertificates) == 0 {
						return fmt.Errorf("no certificates provided by server")
					}
					cert := cs.PeerCertificates[0]
					pubKeyHash := sha256.Sum256(cert.RawSubjectPublicKeyInfo)
					actualHash := fmt.Sprintf("sha256:%x", pubKeyHash)
					if actualHash != envApiHash {
						return fmt.Errorf("public key hash mismatch: expected %s, got %s", envApiHash, actualHash)
					}
					return nil
				},
			},
		},
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
			err := updatePriceCache()
			if err != nil {
				log.Printf("using old cache: %s", err)
			}
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

	err = json.NewDecoder(response.Body).Decode(&pricingResponse)
	if err != nil {
		return &pricingResponse, err
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

	err = json.NewDecoder(response.Body).Decode(&currencyResponse)
	if err != nil {
		return &currencyResponse, err
	}

	return &currencyResponse, nil
}

func updatePriceCache() error {
	priceResult, err := getPrice()
	if err != nil {
		return err
	}
	currencyResult, err := getCurrency()
	if err != nil {
		return err
	}
	priceCache, err := priceResult.toCache(currencyResult.flatten())
	if err != nil {
		return err
	}
	PricingCache = *priceCache
	log.Printf("Updated Price cache on %s", time.Now().String())
	return nil
}

func installItemCache() {
	cacheFile, found := os.LookupEnv("ITEM_CONSTANTS")
	if !found {
		log.Fatalf("ITEM_CONSTANTS is undefined")
	}
	content, err := os.ReadFile(cacheFile)
	if err != nil {
		log.Fatal(err)
	}

	var items = make(map[string]itemConstants)
	err = json.Unmarshal(content, &items)
	if err != nil {
		log.Fatal(err)
	}

	cacheFile, found = os.LookupEnv("UNUSUAL_CONSTANTS")
	if !found {
		log.Fatalf("UNUSUAL_CONSTANTS is undefined")
	}

	content, err = os.ReadFile(cacheFile)
	if err != nil {
		log.Fatal(err)
	}

	var unusuals = make(map[string]string)
	err = json.Unmarshal(content, &unusuals)
	if err != nil {
		log.Fatal(err)
	}

	itemCache = items
	unusualCache = unusuals
	defindexCache = make(map[string]uint32)

	for s, constants := range items {
		defindex, err := strconv.Atoi(s)
		if err != nil {
			continue
		}
		if _, ok := defindexCache[constants.MarketHashName]; ok {
			continue
		}
		defindexCache[constants.MarketHashName] = uint32(defindex)
		defindexCache[fmt.Sprintf("The %s", constants.MarketHashName)] = uint32(defindex)
	}
}

func GetDefindex(itemName string) uint32 {
	if defindex, ok := defindexCache[itemName]; ok {
		return defindex
	}

	trimmed := strings.TrimPrefix(itemName, "Unusual ")
	trimmed = strings.TrimPrefix(trimmed, "Strange ")

	return defindexCache[trimmed]
}

func GetMarketHashName(defindex uint32) string {
	defindexString := strconv.Itoa(int(defindex))

	return itemCache[defindexString].MarketHashName
}
