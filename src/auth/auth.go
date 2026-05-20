package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	c.Status(http.StatusOK)
}

func Login(c *gin.Context) {
	c.Status(http.StatusOK)
}

func Logout(c *gin.Context) {
	c.Status(http.StatusOK)
}

func Me(c *gin.Context) {
	c.Status(http.StatusOK)
}
