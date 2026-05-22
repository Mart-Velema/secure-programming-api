package mfa

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterTOPT(c *gin.Context) {
	c.Status(http.StatusOK)
}

func VerifyTOPT(c *gin.Context) {
	c.Status(http.StatusOK)
}

func ResetTOPT(c *gin.Context) {
	c.Status(http.StatusOK)
}
