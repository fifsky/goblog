package handler

import (
	"github.com/fifsky/goblog/core"
	"github.com/fifsky/goblog/helpers"
	"github.com/fifsky/goblog/models"
	"github.com/gin-gonic/gin"
	"github.com/ilibs/gosql"
	"github.com/verystar/logger"
)

var AdminCommentGet core.HandlerFunc = func(c *core.Context) core.Response {

	h := gin.H{}
	num := 10
	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	comments, err := models.CommentList(page, num)
	h["Comments"] = comments

	total, err := gosql.Model(&models.Comments{}).Count()
	h["Pager"] = c.Pagination(total, num, page)

	if err != nil {
		return c.ErrorMessage(err)
	}

	return c.HTML("admin/comments", h)
}

//var AdminCommentPost core.HandlerFunc = func(c *core.Context) core.Response {
//	return c.Success(nil)
//}

var AdminCommentDelete core.HandlerFunc = func(c *core.Context) core.Response {
	id := helpers.StrTo(c.Query("id")).MustInt()
	if _, err := gosql.Model(&models.Comments{Id: id}).Delete(); err != nil {
		logger.Error(err)
		return c.Fail(201, "删除失败")
	}
	return c.Redirect("/admin/comments")
}
