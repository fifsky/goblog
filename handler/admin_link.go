package handler

import (
	"github.com/fifsky/goblog/helpers"
	"github.com/fifsky/goblog/models"
	"github.com/ilibs/gosql"
	"net/http"
	"github.com/ilibs/logger"
	"github.com/fifsky/goblog/core"
)

var AdminLinkGet core.HandlerFunc = func(c *core.Context) core.Response {
	h := defaultH(c.Context)

	id, _ := helpers.StrTo(c.Query("id")).Int()
	if id > 0 {
		link := &models.Links{Id: id}
		gosql.Model(link).Get()
		h["Link"] = link
	}

	num := 10

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	links, err := models.LinkGetList(page, num)
	h["Links"] = links

	total, err := gosql.Model(&models.Links{}).Count()
	h["Pager"] = c.Pagination(total, num, page)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return nil
	}

	return c.HTML("admin/links", h)
}

var AdminLinkPost core.HandlerFunc = func(c *core.Context) core.Response {
	links := &models.Links{}
	if err := c.Bind(links); err != nil {
		return c.Fail(201, "参数错误:"+err.Error())
	}

	if links.Name == "" {
		return c.Fail(201, "连接名称不能为空")
	}

	if links.Url == "" {
		return c.Fail(201, "连接地址不能为空")
	}

	if links.Id > 0 {
		if _, err := gosql.Model(links).Update(); err != nil {
			logger.Error(err)
			return c.Fail(201, "更新失败")
		}
	} else {
		if _, err := gosql.Model(links).Create(); err != nil {
			logger.Error(err)
			return c.Fail(201, "创建失败")
		}
	}

	return c.Success(nil)
}

var AdminLinkDelete core.HandlerFunc = func(c *core.Context) core.Response {
	id := helpers.StrTo(c.Query("id")).MustInt()

	if _, err := gosql.Model(&models.Links{Id: id}).Delete(); err != nil {
		logger.Error(err)
		return c.Fail(201, "删除失败")
	}
	return c.Redirect("/admin/links")
}
