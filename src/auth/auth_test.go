package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	_ "guineatrade.nhlstenden.com/src/1nit"
	"guineatrade.nhlstenden.com/src/database"
)

func TestRegister(t *testing.T) {
	router := gin.Default()
	router.POST("/login", Register)

	type test struct {
		User       *registerUser
		StatusCode int
		Result     string
	}

	tests := []test{
		{
			&registerUser{
				Email:          "john@doe.com",
				Name:           "John Doe",
				Password:       "1234",
				PasswordVerify: "1234",
			},
			http.StatusCreated,
			"",
		},
		{
			&registerUser{
				Email:          "johndoe.com",
				Name:           "John Doe",
				Password:       "1234",
				PasswordVerify: "1234",
			},
			http.StatusUnprocessableEntity,
			"{\"error\":\"Email is not valid\"}",
		},
		{
			&registerUser{
				Email:          "john@doe.com",
				Name:           "John Doe",
				Password:       "1233",
				PasswordVerify: "1234",
			},
			http.StatusUnprocessableEntity,
			"{\"error\":\"Passwords do not match\"}",
		},
		{
			&registerUser{
				Email:          "john@doe.com",
				Name:           "John Doe",
				Password:       "1234",
				PasswordVerify: "1234",
			},
			http.StatusConflict,
			"{\"error\":\"User already exists\"}",
		},
	}

	for _, t2 := range tests {
		userJson, _ := json.Marshal(t2.User)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/login", bytes.NewReader(userJson))

		router.ServeHTTP(w, req)

		assert.Equal(t, t2.StatusCode, w.Code)
		assert.Equal(t, t2.Result, w.Body.String())
	}
}

func TestLogin(t *testing.T) {
	password := "abcd"
	user := database.CreateRandomUser()
	user.Password = password
	database.GetInstance().Select("password").Save(user)
	database.GetInstance().First(user)

	router := gin.Default()
	router.POST("/login", Login)

	type test struct {
		User       *registerUser
		StatusCode int
		Result     string
	}

	tests := []test{
		{
			&registerUser{
				Email:    "johan@doe.com",
				Password: "1234",
			},
			http.StatusNotFound,
			"{\"error\":\"User does not exist\"}",
		},
		{
			&registerUser{
				Email:    user.Email,
				Password: password,
			},
			http.StatusOK,
			"jwt",
		},
	}

	for _, t2 := range tests {
		userJson, _ := json.Marshal(t2.User)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/login", bytes.NewReader(userJson))

		router.ServeHTTP(w, req)

		assert.Equal(t, t2.StatusCode, w.Code)
		assert.Equal(t, strings.Contains(w.Body.String(), t2.Result), true)
	}
}

func TestIsEmailValid(t *testing.T) {
	type test struct {
		Email   string
		IsValid bool
	}

	tests := []test{
		{"user@example.com", true},
		{"first.last@domain.co.uk", true},
		{"username+tag@group.mailprovider.org", true},
		{"firstname.lastname@domain.com", true},
		{"email@localhost", false},
		{"user-name@sub.domain.example.com", true},
		{"1234567890@domain.com", true},
		{"username_123@domain.co.in", true},
		{"用户@example.com", false},
		{"пользователь@домен.рф", false},
		{"plainaddress", false},
		{"#<http://example.com/>", false},
		{"<EMAIL>", false},
		{"Joe Smith <email@example.com>", false},
		{"email.example.com", false},
		{"@domain.com", false},
		{"user@.com", false},
		{"user@domain.toolongtld", false},
		{"username@-domain.com", true},
		{"username@domain-.com", true},
		{"username@domain..com", true},
	}

	for _, t2 := range tests {
		result := isEmailValid(t2.Email)
		assert.Equal(t, result, t2.IsValid)
	}
}
