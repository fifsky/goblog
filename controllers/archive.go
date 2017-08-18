package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/fifsky/goblog/models"
	"github.com/fifsky/goblog/helpers"
)

func ArchiveGet(c *gin.Context) {
	id, _ := helpers.StrTo(c.Param("id")).Uint()
	postModel := &models.Posts{Id: id}
	post, err := postModel.Get()
	if err == nil {
		c.HTML(http.StatusOK, "index/article", gin.H{
			"post": post,
		})
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
