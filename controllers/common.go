package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SESSION_KEY          = "UserID"       // session key
	CONTEXT_USER_KEY     = "User"         // context user key
)

func Handle404(c *gin.Context) {
	HandleMessage(c, "Sorry,I lost myself!")
}

func HandleMessage(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{
		"message": message,
	})
}