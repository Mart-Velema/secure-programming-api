package inventory

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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"guineatrade.nhlstenden.com/src/auth/middleware"
	"guineatrade.nhlstenden.com/src/items"
)

var steamClient *http.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s\n", err)
	}

	envSteamHash, steamHashExists := os.LookupEnv("STEAM_API_HASH")

	if !steamHashExists {
		log.Fatal("STEAM_API_HASH is unset")
	}

	steamClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    30,
			MaxConnsPerHost: 10,
			IdleConnTimeout: 30 * time.Second,
			TLSClientConfig: &tls.Config{
				ServerName:         "steamcommunity.com",
				InsecureSkipVerify: false,
				VerifyConnection: func(cs tls.ConnectionState) error {
					if len(cs.PeerCertificates) == 0 {
						return fmt.Errorf("no certificates provided by server")
					}

					cert := cs.PeerCertificates[0]

					pubKeyHash := sha256.Sum256(
						cert.RawSubjectPublicKeyInfo,
					)

					actualHash := fmt.Sprintf(
						"sha256:%x",
						pubKeyHash,
					)

					if actualHash != envSteamHash {
						return fmt.Errorf(
							"public key hash mismatch: expected %s, got %s",
							envSteamHash,
							actualHash,
						)
					}

					return nil
				},
			},
		},
	}

	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := tls.Dial(
		"tcp",
		"steamcommunity.com:443",
		&tls.Config{
			RootCAs:    rootCAs,
			ServerName: "steamcommunity.com",
		},
	)

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
		DNSName:       "steamcommunity.com",
		Intermediates: intermediates,
	}

	_, err = chain[0].Verify(opts)
	if err != nil {
		log.Fatal(err)
	}
}

func GetUserInventory(steamID uint64) (*items.InventoryResponse, error) {
	url := fmt.Sprintf(
		"https://steamcommunity.com/inventory/%d/440/2?count=2000&l=english",
		steamID,
	)

	response, err := steamClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"steam returned status %d",
			response.StatusCode,
		)
	}

	var inventory items.InventoryResponse

	err = json.NewDecoder(response.Body).Decode(&inventory)
	if err != nil {
		return nil, err
	}

	return &inventory, nil
}

func GetInventory(c *gin.Context) {
	user, err := middleware.ExtractTokenUser(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Token expired"})
		return
	}

	inventory, err := GetUserInventory(user.SteamId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get inventory data"})
		return
	}
	c.JSON(200, inventory.ToItem().Assets)
}
