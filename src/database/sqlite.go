package database

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/cgholdings/go-common/database/encryption"
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
	Balance      int64          `json:"-" gorm:"default:0"`
	SteamId      uint64         `json:"steamId"`
	TradeUrl     string         `json:"tradeUrl" encrypt:"true" gorm:"unique"`
	TotpSecret   string         `encrypt:"true"`
	RecoveryCode string         `hash:"true"`
	Trades       []Trade        `gorm:"foreignKey:UserID"`
	Token        []RefreshToken `gorm:"foreignKey:UserID"`
}

func (u User) HasMFAEnabled() bool {
	return len(u.RecoveryCode) != 0
}

func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Balance  int64  `json:"balance"`
		HasMFA   bool   `json:"mfaEnabled"`
		SteamId  uint64 `json:"steamId"`
		TradeUrl string `json:"tradeUrl"`
	}{
		Name:     u.Name,
		Email:    u.Email,
		Balance:  u.Balance,
		HasMFA:   u.HasMFAEnabled(),
		SteamId:  u.SteamId,
		TradeUrl: u.TradeUrl,
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

func createEncryptor() {
	if instanceEncryptor != nil {
		return
	}

	encryptor, err := encryption.NewEncryptorFromConfig(encryption.DefaultConfig())
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
