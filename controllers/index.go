package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/models"
	"net/http"
	"fmt"
)

func IndexGet(c *gin.Context) {
	post := &models.Posts{Id:5}
	posts, err := post.GetList()

	fmt.Println(err)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"posts": posts,
		})
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func AdminIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/index", gin.H{
	})
}
