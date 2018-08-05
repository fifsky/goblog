package handler

import (
	"net/http"

	"github.com/fifsky/goblog/helpers"
	"github.com/fifsky/goblog/models"
	"github.com/ilibs/gosql"
	"github.com/ilibs/logger"
	"github.com/fifsky/goblog/core"
)

var AdminUsersGet core.HandlerFunc = func(c *core.Context) core.Response {
	h := defaultH(c.Context)
	num := 10

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	users, err := models.UserGetList(page, num)

	h["Users"] = users

	total, err := gosql.Model(&models.Users{}).Count()
	h["Pager"] = c.Pagination(total, num, page)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return nil
	}
	return c.HTML("admin/users", h)
}

var AdminUserGet core.HandlerFunc = func(c *core.Context) core.Response {
	h := defaultH(c.Context)
	id := helpers.StrTo(c.Query("id")).MustInt()
	if id > 0 {
		user := &models.Users{Id: id}
		err := gosql.Model(user).Get()
		if err != nil {
			return HandleMessage(c, "用户不存在", "您访问的用户不存在或已经删除！")
		}
		h["User"] = user
	}

	return c.HTML("admin/post_user", h)
}

var AdminUserPost core.HandlerFunc = func(c *core.Context) core.Response {
	users := &models.Users{}
	if err := c.Bind(users); err != nil {
		return c.Fail(201, "参数错误:"+err.Error())
	}

	if users.Name == "" {
		return c.Fail(201, "用户名不能为空")
	}

	if users.Password == "" {
		return c.Fail(201, "密码不能为空")
	}

	if users.NickName == "" {
		return c.Fail(201, "昵称不能为空")
	}

	users.Password = helpers.Md5(users.Password)

	if users.Id > 0 {
		if _, err := gosql.Model(users).Update(); err != nil {
			logger.Error(err)
			return c.Fail(201, "更新失败")
		}
	} else {
		if _, err := gosql.Model(users).Create(); err != nil {
			logger.Error(err)
			return c.Fail(201, "创建失败")
		}
	}

	return c.Success(nil)
}

var AdminUserStatus core.HandlerFunc = func(c *core.Context) core.Response {
	id := helpers.StrTo(c.Query("id")).MustInt()
	status := helpers.StrTo(c.Query("status")).MustInt()

	if _, err := gosql.Model(&models.Users{}).Where("id = ? and status = ?", id, status).Delete(); err != nil {
		logger.Error(err)
		return c.Fail(201, "删除失败")
	}
	return c.Redirect("/admin/users")
}
