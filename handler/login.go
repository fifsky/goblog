package handler

import (
	"github.com/fifsky/goblog/models"
	"github.com/ilibs/gosql"
	"github.com/fifsky/goblog/helpers"
	"github.com/fifsky/goblog/core"
)

var LoginGet core.HandlerFunc = func(c *core.Context) core.Response {
	h := defaultH(c.Context)
	if h["LoginUser"] != nil {
		return c.Redirect("/admin/index")
	}
	return c.HTML("admin/login", h)
}

var LogoutGet core.HandlerFunc = func(c *core.Context) core.Response {
	c.Session().Delete("UserId")
	c.Session().Save()
	return c.Redirect("/")
}

var LoginPost core.HandlerFunc = func(c *core.Context) core.Response {
	user_name := c.PostForm("user_name")
	password := c.PostForm("user_pass")

	if user_name == "" || password == "" {
		return c.Fail(201, "用户名或密码不能为空")
	}

	user := &models.Users{Name: user_name}
	err := gosql.Model(user).Get()

	if err != nil {
		return c.Fail(201, "用户名或者密码错误")
	}

	if user.Password != helpers.Md5(password) {
		return c.Fail(201, "用户名或者密码错误")
	}

	c.Session().Set("UserId", user.Id)
	c.Session().Save()

	return c.Success(nil)
}
