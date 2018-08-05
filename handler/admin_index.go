package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/models"
	"github.com/ilibs/gosql"
	"github.com/fifsky/goblog/core"
)

func AdminIndex(c *gin.Context) {
	h := defaultH(c)
	h["Options"] = c.GetStringMapString("options")
	c.HTML(http.StatusOK, "admin/index", h)
}

func AdminIndexPost(c *gin.Context) {
	c.Request.ParseForm()
	options := c.Request.PostForm

	for k, v := range options {
		gosql.Model(&models.Options{
			OptionValue: v[0],
		}).Where("option_key = ?", k).Update()
	}

	o, err := models.GetOptions()

	if err == nil {
		core.Global.Store("options", o)
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"message":    "保存成功",
		"options":    options,
	})
}
