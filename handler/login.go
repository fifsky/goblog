package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/fifsky/goblog/models"
	"github.com/ilibs/gosql"
	"github.com/fifsky/goblog/helpers"
	"github.com/ilibs/sessions"
)

func LoginGet(c *gin.Context) {
	h := defaultH(c)
	if h["LoginUser"] != nil {
		c.Redirect(http.StatusFound, "/admin/index")
		return
	}

	c.HTML(http.StatusOK, "admin/login", h)
}

func LogoutGet(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("UserId")
	session.Save()
	c.Redirect(http.StatusFound, "/")
}

func LoginPost(c *gin.Context) {
	session := sessions.Default(c)

	user_name := c.PostForm("user_name")
	password := c.PostForm("user_pass")

	if user_name == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "用户名或密码不能为空",
		})
		return
	}

	user := &models.Users{Name: user_name}
	err := gosql.Model(user).Get()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "用户名或者密码错误",
		})
		return
	}

	if user.Password != helpers.Md5(password) {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "密码错误",
		})
		return
	}

	session.Set("UserId", user.Id)
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"message":    "ok",
	})
}
