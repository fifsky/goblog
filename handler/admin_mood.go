package handler

import (
	"github.com/fifsky/goblog/core"
	"github.com/fifsky/goblog/helpers"
	"github.com/fifsky/goblog/models"
	"github.com/gin-gonic/gin"
	"github.com/ilibs/gosql"
	"github.com/verystar/logger"
)

var AdminMoodGet core.HandlerFunc = func(c *core.Context) core.Response {
	h := gin.H{}

	id := helpers.StrTo(c.Query("id")).MustInt()
	if id > 0 {
		mood := &models.Moods{Id: id}
		gosql.Model(mood).Get()
		h["Mood"] = mood
	}

	num := 10

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	moods, err := models.MoodGetList(page, num)
	h["Moods"] = moods

	total, err := gosql.Model(&models.Moods{}).Count()
	h["Pager"] = c.Pagination(total, num, page)

	if err != nil {
		return c.ErrorMessage(err)
	}
	return c.HTML("admin/moods", h)
}

var AdminMoodPost core.HandlerFunc = func(c *core.Context) core.Response {
	moods := &models.Moods{}
	if err := c.Bind(moods); err != nil {
		return c.Fail(201, "参数错误:"+err.Error())
	}

	moods.UserId = c.Session().Get("UserId").(int)

	if moods.Content == "" {
		return c.Fail(201, "内容不能为空")
	}

	if moods.Id > 0 {
		if _, err := gosql.Model(moods).Update(); err != nil {
			logger.Error(err)
			return c.Fail(201, "更新失败")
		}
	} else {
		if _, err := gosql.Model(moods).Create(); err != nil {
			logger.Error(err)
			return c.Fail(201, "发表失败")
		}
	}
	return c.Success(nil)
}

var AdminMoodDelete core.HandlerFunc = func(c *core.Context) core.Response {
	id := helpers.StrTo(c.Query("id")).MustInt()

	mood := &models.Moods{Id: id}
	if _, err := gosql.Model(mood).Delete(); err != nil {
		logger.Error(err)
		return c.Fail(201, "删除失败")
	}
	return c.Redirect(c.Request.Referer())
}
