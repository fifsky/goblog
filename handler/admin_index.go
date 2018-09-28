package handler

import (
	"github.com/fifsky/goblog/core"
	"github.com/fifsky/goblog/models"
	"github.com/ilibs/gosql"
)

var AdminIndex core.HandlerFunc = func(c *core.Context) core.Response {
	return c.HTML("admin/index")
}

var AdminIndexPost core.HandlerFunc = func(c *core.Context) core.Response {
	c.Request.ParseForm()
	options := c.Request.PostForm

	for k, v := range options {
		gosql.Model(&models.Options{
			OptionValue: v[0],
		}).Where("option_key = ?", k).Update()
	}

	return c.Success(options)
}
