package mfa

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendSMS(c *gin.Context) {
	c.Status(http.StatusOK)
}

func VerifySMS(c *gin.Context) {
	c.Status(http.StatusOK)
}
