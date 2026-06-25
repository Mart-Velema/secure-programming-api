package database

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/assert/v2"
	_ "guineatrade.nhlstenden.com/src/1nit"
)

func TestGetInstance(t *testing.T) {
	db := GetInstance()

	assert.Equal(t, db, GetInstance())
}

func TestCreateRandomUser(t *testing.T) {
	user := CreateRandomUser()

	assert.NotEqual(t, user, CreateRandomUser())
}

func TestUser_HasMFAEnabled(t *testing.T) {
	user := CreateRandomUser()

	assert.Equal(t, user.HasMFAEnabled(), true)

	user.RecoveryCode = ""
	user.TotpSecret = ""

	GetInstance().Select("recovery_code", "totp_secret").Save(user)

	assert.Equal(t, user.HasMFAEnabled(), false)
}

func TestUser_MarshalJSON(t *testing.T) {
	user := User{
		Name:     "john@doe.com",
		Email:    "john@doe.com",
		Password: "1234",
	}

	GetInstance().Save(&user)
	GetInstance().First(&user)

	userJson, _ := json.Marshal(&user)
	assert.Equal(t, string(userJson), "{\"name\":\"john@doe.com\",\"email\":\"john@doe.com\",\"balance\":0,\"mfaEnabled\":false,\"steamId\":0,\"tradeUrl\":\"\"}")
}
