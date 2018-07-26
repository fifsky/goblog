package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/helpers"
	"github.com/fifsky/goblog/models"
	"github.com/ilibs/gosql"
	"github.com/fifsky/goblog/helpers/pagination"
	"github.com/ilibs/logger"
)

func AdminUsersGet(c *gin.Context) {
	h := defaultH(c)
	num := 10

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	users, err := models.UserGetList(page, num)

	h["Users"] = users

	total, err := gosql.Model(&models.Users{}).Count()
	pager := pagination.New(int(total), num, page, 3)
	h["Pager"] = pager

	if err == nil {
		c.HTML(http.StatusOK, "admin/users", h)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func AdminUserGet(c *gin.Context) {
	h := defaultH(c)
	id := helpers.StrTo(c.Query("id")).MustInt()
	if id > 0 {
		user := &models.Users{Id: id}
		err := gosql.Model(user).Get()
		if err != nil {
			HandleMessage(c, "用户不存在", "您访问的用户不存在或已经删除！")
			return
		}
		h["User"] = user
	}

	c.HTML(http.StatusOK, "admin/post_user", h)
}

func AdminUserPost(c *gin.Context) {
	users := &models.Users{}
	if err := c.Bind(users); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "参数错误:" + err.Error(),
		})
		return
	}

	if users.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "用户名不能为空",
		})
		return
	}

	if users.Password == "" {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "密码不能为空",
		})
		return
	}

	if users.NickName == "" {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "昵称不能为空",
		})
		return
	}

	users.Password = helpers.Md5(users.Password)

	if users.Id > 0 {
		if _, err := gosql.Model(users).Update(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "更新失败",
			})
			logger.Error(err)
			return
		}
	} else {
		if _, err := gosql.Model(users).Create(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "创建失败",
			})
			logger.Error(err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"message":    "ok",
	})
}

func AdminUserStatus(c *gin.Context) {
	id := helpers.StrTo(c.Query("id")).MustInt()
	status := helpers.StrTo(c.Query("status")).MustInt()

	if _, err := gosql.Model(&models.Users{}).Where("id = ? and status = ?", id, status).Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "删除失败",
		})
		logger.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/users")
}
