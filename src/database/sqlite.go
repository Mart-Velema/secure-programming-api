package database

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/cgholdings/go-common/database/encryption"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	lockSqlite = &sync.Mutex{}
	instanceDB *gorm.DB
)

type User struct {
	gorm.Model   `json:"-"`
	Name         string         `json:"name" gorm:"unique"`
	Email        string         `json:"email" encrypt:"true"`
	EmailHash    string         `json:"-" hash:"Email" gorm:"unique"`
	Password     string         `json:"password" hash:"Password"`
	Balance      int64          `json:"-" gorm:"default:0"`
	SteamId      uint64         `json:"steamId"`
	TradeUrl     string         `json:"tradeUrl" encrypt:"true"`
	TotpSecret   string         `encrypt:"true"`
	RecoveryCode string         `hash:"RecoveryCode"`
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
	gorm.Model       `json:"-"`
	UserID           uint
	Cost             int64
	TradeAction      TradeAction
	TradeStatus      TradeStatus
	SteamTradeId     string
	StripePaymentUrl string
	Assets           []Asset
}

type Asset struct {
	gorm.Model     `json:"-"`
	TradeID        uint
	Trade          Trade
	TradeDirection TradeAction
	AssetId        string
}

type TradeAction int

const (
	BUY TradeAction = iota
	SELL
)

type TradeStatus int

const (
	PAYMENT_IN_PROGRESS TradeStatus = iota
	TRADE_IN_PROGRESS
	COMPLETED
	CANCELLED
)

func createEncryptor() *encryption.Encryptor {
	encryptor, err := encryption.NewEncryptorFromConfig(encryption.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	return encryptor
}

func createDB() {
	if instanceDB != nil {
		return
	}
	db, err := gorm.Open(sqlite.Open(os.Getenv("SQLITE_FILE_LOCATION")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	encryptor := createEncryptor()
	err = db.Use(encryption.NewPlugin(encryptor))
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&User{}, &Trade{}, &Asset{}, &RefreshToken{})
	if err != nil {
		log.Fatal(err)
	}

	instanceDB = db
}

func CreateRandomUser() *User {
	password := uuid.New()
	mfaCode := rand.Text()
	email := fmt.Sprintf("%s@%s.com", uuid.New(), uuid.New())
	name := rand.Text()

	user := User{
		Email:        email,
		Name:         name,
		Password:     password.String(),
		TotpSecret:   mfaCode,
		RecoveryCode: mfaCode,
	}

	GetInstance().Save(&user)
	GetInstance().First(&user)

	return &user
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
