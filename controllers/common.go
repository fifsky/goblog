package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handle404(c *gin.Context) {
	HandleMessage(c, "未找到(404 Not Found)", "抱歉，您浏览的页面未找到。")
}

func HandleMessage(c *gin.Context, title string, message string) {
	c.HTML(http.StatusNotFound, "error/message", gin.H{
		"Title":   title,
		"Message": message,
	})
}