package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/cgholdings/go-common/database/encryption"
	"golang.org/x/crypto/argon2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	lockSqlite        = &sync.Mutex{}
	lockEncryptor     = &sync.Mutex{}
	instanceDB        *gorm.DB
	instanceEncryptor *encryption.Encryptor
)

type User struct {
	gorm.Model   `json:"-"`
	Name         string         `json:"name" gorm:"unique"`
	Email        string         `json:"email" encrypt:"true"`
	EmailHash    string         `json:"-" hash:"Email" gorm:"unique"`
	Password     string         `json:"password" hash:"Password"`
	PhoneNumber  string         `json:"tel" encrypt:"true"`
	NumberHash   string         `json:"-" hash:"PhoneNumber" gorm:"unique"`
	Balance      int64          `json:"-" gorm:"default:0"`
	TotpSecret   string         `encrypt:"true"`
	RecoveryCode string         `hash:"true"`
	Trades       []Trade        `gorm:"foreignKey:UserID"`
	Token        []RefreshToken `gorm:"foreignKey:UserID"`
}

func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"tel"`
		Balance     int64  `json:"balance"`
	}{
		Name:        u.Name,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Balance:     u.Balance,
	})
}

type RefreshToken struct {
	UserID    uint   `gorm:"index"`
	User      User   `gorm:"foreignKey:UserID"`
	Token     string `encrypt:"true"`
	TokenHash string `hash:"Token" gorm:"uniqueIndex"`
	Nonce     string
	ExpiresOn time.Time
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

func createEncryptor() {
	if instanceEncryptor != nil {
		return
	}

	config := encryption.DefaultConfig()

	if key, exists := os.LookupEnv("ENCRYPTION_PASSCODE"); exists {
		derived := deriveKey(key)
		fmt.Printf("Derived key: %x\n", derived)
		config.Key = derived
	}

	encryptor, err := encryption.NewEncryptorFromConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	instanceEncryptor = encryptor
}

func GetEncryptor() *encryption.Encryptor {
	if instanceEncryptor != nil {
		return instanceEncryptor
	}
	lockEncryptor.Lock()
	defer lockEncryptor.Unlock()
	createEncryptor()

	return instanceEncryptor
}

func createDB() {
	if instanceDB != nil {
		return
	}
	db, err := gorm.Open(sqlite.Open(os.Getenv("SQLITE_FILE_LOCATION")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	encryptor := GetEncryptor()
	err = db.Use(encryption.NewPlugin(encryptor))
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&User{}, &Trade{}, &TradeItem{}, &RefreshToken{})
	if err != nil {
		log.Fatal(err)
	}

	instanceDB = db
}

func GetInstance() *gorm.DB {
	if instanceDB != nil {
		return instanceDB
	}
	lockSqlite.Lock()
	defer lockSqlite.Unlock()
	createDB()

	return instanceDB
}
