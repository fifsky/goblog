package handler

import (
	"github.com/fifsky/goblog/models"
	"github.com/ilibs/gosql"
	"github.com/fifsky/goblog/core"
)

var AdminIndex core.HandlerFunc = func(c *core.Context) core.Response {
	h := defaultH(c.Context)
	h["Options"] = c.GetStringMapString("options")
	return c.HTML("admin/index", h)
}

var AdminIndexPost core.HandlerFunc = func(c *core.Context) core.Response {
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

	return c.Success(options)
}
