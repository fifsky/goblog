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

func AdminMoodGet(c *gin.Context) {
	h := defaultH(c)

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
	pager := pagination.New(int(total), num, page, 3)
	h["Pager"] = pager

	if err == nil {
		c.HTML(http.StatusOK, "admin/moods", h)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func AdminMoodPost(c *gin.Context) {
	moods := &models.Moods{}
	if err := c.Bind(moods); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "参数错误:" + err.Error(),
		})
		return
	}

	if user, exists := c.Get("LoginUser"); exists {
		moods.UserId = user.(*models.Users).Id
	}

	if moods.Content == "" {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "内容不能为空",
		})
		return
	}

	if moods.Id > 0 {
		if _, err := gosql.Model(moods).Update(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "更新失败",
			})
			logger.Error(err)
			return
		}
	} else {
		if _, err := gosql.Model(moods).Create(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "发表失败",
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

func AdminMoodDelete(c *gin.Context) {
	id := helpers.StrTo(c.Query("id")).MustInt()

	mood := &models.Moods{Id: id}
	if _, err := gosql.Model(mood).Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "删除失败",
		})
		logger.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/moods")
}
