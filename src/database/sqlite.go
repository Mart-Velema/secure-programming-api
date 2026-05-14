package database

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"

	"github.com/cgholdings/go-common/database/encryption"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model  `json:"-"`
	Name        string `json:"name"`
	Email       string `json:"email" encrypt:"true"`
	PhoneNumber string `json:"phone_number" encrypt:"true"`
}

func CreateDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("guineatrade.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	config := encryption.DefaultConfig()

	if key, exists := os.LookupEnv("ENCRYPTION_PASSCODE"); exists {
		hasher := sha256.New()
		hasher.Write([]byte(key))
		config.Key = hasher.Sum(nil)
	}

	encryptor, err := encryption.NewEncryptorFromConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Use(encryption.NewPlugin(encryptor))
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}
	user := &User{
		Name:        "John Doe",
		Email:       "john.doe@example.com",
		PhoneNumber: "+1-555-1234-5678",
	}
	db.Create(user)

	var users []User
	if result := db.Find(&users); result.Error != nil {
		log.Fatal(err)
	}

	for idx, user := range users {
		fmt.Printf("User %d:\n\tName: %s\n\tEmail: %s\n\tPhone: %s\n", idx, user.Name, user.Email, user.PhoneNumber)
	}

	{
		var user User
		user.ID = 1
		db.Find(&user)
		fmt.Printf("\n\n User from database %d:\n\tName: %s\n\tEmail: %s\n\tPhone: %s\n", user.ID, user.Name, user.Email, user.PhoneNumber)
	}

	return db
}
