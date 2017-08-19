package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/models"
	"net/http"
	"github.com/fifsky/goblog/helpers"
	"github.com/fifsky/goblog/helpers/pagination"
)

func IndexGet(c *gin.Context) {

	num := 10

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	postModel := new(models.Posts)
	posts, err := postModel.GetList(page, num)

	h := defaultH(c)
	h["Posts"] = posts

	total, err := postModel.Count()
	pager := pagination.New(int(total), num, page, 3)
	h["Pager"] = pager

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
