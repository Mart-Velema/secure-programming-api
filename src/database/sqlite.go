package database

import (
	"log"
	"os"

	"github.com/cgholdings/go-common/database/encryption"
	"golang.org/x/crypto/argon2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model  `json:"-"`
	Name        string `json:"name"`
	Email       string `json:"email" encrypt:"true"`
	EmailHash   string `hash:"Email"`
	Password    string `json:"password" hash:"Password"`
	PhoneNumber string `json:"phone_number" encrypt:"true"`
	NumberHash  string `hash:"PhoneNumber"`
	Balance     int64
	Trades      []Trade `gorm:"foreignKey:UserID"`
}

type Trade struct {
	gorm.Model  `json:"-"`
	UserID      uint
	Cost        int64
	SoldItems   []TradeItem `gorm:"foreignKey:TradeID"`
	BoughtItems []TradeItem `gorm:"foreignKey:TradeID"`
}

type TradeItem struct {
	gorm.Model `json:"-"`
	TradeID    uint
	ItemID     uint
	Quantity   uint
}

func deriveKey(passcode string) []byte {
	bytesKey := []byte(passcode)
	mid := len(bytesKey) / 2
	salt := bytesKey[mid : mid+16]

	return argon2.IDKey(
		bytesKey,
		salt,
		3,
		64*1024,
		4,
		32,
	)
}

func CreateDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("guineatrade.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	config := encryption.DefaultConfig()

	if key, exists := os.LookupEnv("ENCRYPTION_PASSCODE"); exists {
		config.Key = deriveKey(key)
	}

	encryptor, err := encryption.NewEncryptorFromConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Use(encryption.NewPlugin(encryptor))
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&User{}, &Trade{}, &TradeItem{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
