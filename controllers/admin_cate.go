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

func AdminCateGet(c *gin.Context) {
	h := defaultH(c)

	id := helpers.StrTo(c.Query("id")).MustInt()
	if id > 0 {
		cate := &models.Cates{Id: id}
		gosql.Model(cate).Get()
		h["Cate"] = cate
	}

	num := 10

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	cates, err := models.CateArtivleCountGetList(page, num)
	h["Cates"] = cates

	total, err := gosql.Model(&models.Cates{}).Count()
	pager := pagination.New(int(total), num, page, 3)
	h["Pager"] = pager

	if err == nil {
		c.HTML(http.StatusOK, "admin/cates", h)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func AdminCatePost(c *gin.Context) {
	cates := &models.Cates{}
	if err := c.Bind(cates); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "参数错误:" + err.Error(),
		})
		return
	}

	if cates.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "分类名不能为空",
		})
		return
	}

	if cates.Domain == "" {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "分类缩略名不能为空",
		})
		return
	}

	if cates.Id > 0 {
		if _, err := gosql.Model(cates).Update(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "更新失败",
			})
			logger.Error(err)
			return
		}
	} else {
		if _, err := gosql.Model(cates).Create(); err != nil {
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

func AdminCateDelete(c *gin.Context) {
	id := helpers.StrTo(c.Query("id")).MustInt()
	total, _ := gosql.Model(&models.Posts{}).Where("cate_id = ?", id).Count()

	if total > 0 {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "该分类下面还有文章，不能删除",
		})
		return
	}

	if _, err := gosql.Model(&models.Cates{Id: id}).Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "删除失败",
		})
		logger.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/cates")
}
