package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/helpers"
	"github.com/fifsky/goblog/models"
	"github.com/ilibs/gosql"
	"time"
	"net/http"
	"github.com/gin-gonic/gin/binding"
	"github.com/fifsky/goblog/helpers/pagination"
	"github.com/ilibs/logger"
)

func AdminRemindGet(c *gin.Context) {
	h := defaultH(c)

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
	pager := pagination.New(int(total), num, page, 3)
	h["Pager"] = pager

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

	if err == nil {
		c.HTML(http.StatusOK, "admin/remind", h)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func AdminRemindPost(c *gin.Context) {
	reminds := &models.Reminds{}
	if err := c.ShouldBindWith(reminds, binding.Form); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "参数错误:" + err.Error(),
		})
		return
	}

	if reminds.Content == "" {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "提醒内容不能为空",
		})
		return
	}

	if reminds.Id > 0 {
		if _, err := gosql.Model(reminds).Update(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "更新失败:" + err.Error(),
			})
			logger.Error(err)
			return
		}
	} else {
		if _, err := gosql.Model(reminds).Create(); err != nil {
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

func AdminRemindDelete(c *gin.Context) {
	id := helpers.StrTo(c.Query("id")).MustInt()

	if _, err := gosql.Model(&models.Reminds{Id: id}).Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "删除失败",
		})
		logger.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/remind")
}
