package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/helpers"
	"github.com/fifsky/goblog/models"
	"github.com/ilibs/gosql"
	"net/http"
	"github.com/fifsky/goblog/helpers/pagination"
	"github.com/ilibs/logger"
)

func AdminLinkGet(c *gin.Context) {
	h := defaultH(c)

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
	pager := pagination.New(int(total), num, page, 3)
	h["Pager"] = pager

	if err == nil {
		c.HTML(http.StatusOK, "admin/links", h)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func AdminLinkPost(c *gin.Context) {
	links := &models.Links{}
	if err := c.Bind(links); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "参数错误:" + err.Error(),
		})
		return
	}

	if links.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "连接名称不能为空",
		})
		return
	}

	if links.Url == "" {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "连接地址不能为空",
		})
		return
	}

	if links.Id > 0 {
		if _, err := gosql.Model(links).Update(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "更新失败",
			})
			logger.Error(err)
			return
		}
	} else {
		if _, err := gosql.Model(links).Create(); err != nil {
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

func AdminLinkDelete(c *gin.Context) {
	id := helpers.StrTo(c.Query("id")).MustInt()

	if _, err := gosql.Model(&models.Links{Id: id}).Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "删除失败",
		})
		logger.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/links")
}
