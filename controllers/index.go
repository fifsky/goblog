package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/models"
	"net/http"
	"github.com/fifsky/goblog/helpers"
)

func IndexGet(c *gin.Context) {
	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	postModel := new(models.Posts)
	posts, err := postModel.GetList(page, 10)

	h := defaultH(c)
	h["Posts"] = posts

	if err == nil {
		c.HTML(http.StatusOK, "index/index", h)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func AdminIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/index", gin.H{
	})
}
