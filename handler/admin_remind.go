package handler

import (
	"net/http"
	"time"

	"github.com/fifsky/goblog/core"
	"github.com/fifsky/goblog/helpers"
	"github.com/fifsky/goblog/models"
	"github.com/gin-gonic/gin/binding"
	"github.com/ilibs/gosql"
	"github.com/ilibs/logger"
)

var AdminRemindGet core.HandlerFunc = func(c *core.Context) core.Response {
	h := defaultH(c.Context)

	id := helpers.StrTo(c.Query("id")).MustInt()
	if id > 0 {
		remind := &models.Reminds{Id: id}
		gosql.Model(remind).Get()
		h["Remind"] = remind
	}

	num := 10

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	reminds, err := models.RemindGetList(page, num)

	h["Reminds"] = reminds

	total, err := gosql.Model(&models.Reminds{}).Count()
	h["Pager"] = c.Pagination(total, num, page)

	h["Types"] = map[int]string{
		0: "固定",
		1: "每分钟",
		2: "每小时",
		3: "每天",
		4: "每周",
		5: "每月",
		6: "每年",
	}

	h["Layouts"] = map[int]string{
		0: "2006-01-02 15:04:05",
		1: "",
		2: "",
		3: "15:04:00",
		4: "15:04:00",
		5: "02日15:04:05",
		6: "01月02日15:04:05",
	}

	h["CurrDate"] = time.Now().Format("2006-01-02 15:04:05")

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil
	}
	return c.HTML("admin/remind", h)
}

var AdminRemindPost core.HandlerFunc = func(c *core.Context) core.Response {
	reminds := &models.Reminds{}
	if err := c.ShouldBindWith(reminds, binding.Form); err != nil {
		return c.Fail(201, "参数错误:"+err.Error())
	}

	if reminds.Content == "" {
		return c.Fail(201, "提醒内容不能为空")
	}

	if reminds.Id > 0 {
		if _, err := gosql.Model(reminds).Update(); err != nil {
			logger.Error(err)
			return c.Fail(201, "更新失败:"+err.Error())
		}
	} else {
		if _, err := gosql.Model(reminds).Create(); err != nil {
			logger.Error(err)
			return c.Fail(201, "创建失败")
		}
	}

	return c.Success(nil)
}

var AdminRemindDelete core.HandlerFunc = func(c *core.Context) core.Response {
	id := helpers.StrTo(c.Query("id")).MustInt()

	if _, err := gosql.Model(&models.Reminds{Id: id}).Delete(); err != nil {
		logger.Error(err)
		return c.Fail(201, "删除失败")
	}
	return c.Redirect("/admin/remind")
}
