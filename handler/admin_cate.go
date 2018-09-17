package handler

import (
	"net/http"

	"github.com/fifsky/goblog/core"
	"github.com/fifsky/goblog/helpers"
	"github.com/fifsky/goblog/models"
	"github.com/ilibs/gosql"
	"github.com/ilibs/logger"
)

var AdminCateGet core.HandlerFunc = func(c *core.Context) core.Response {
	h := defaultH(c.Context)

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
	h["Pager"] = c.Pagination(total, num, page)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil
	}

	return c.HTML("admin/cates", h)

}

var AdminCatePost core.HandlerFunc = func(c *core.Context) core.Response {
	cates := &models.Cates{}
	if err := c.Bind(cates); err != nil {
		return c.Fail(201, "参数错误:"+err.Error())
	}

	if cates.Name == "" {
		return c.Fail(201, "分类名不能为空")
	}

	if cates.Domain == "" {
		return c.Fail(201, "分类缩略名不能为空")
	}

	if cates.Id > 0 {
		if _, err := gosql.Model(cates).Update(); err != nil {
			logger.Error(err)
			return c.Fail(201, "更新失败")
		}
	} else {
		if _, err := gosql.Model(cates).Create(); err != nil {
			logger.Error(err)
			return c.Fail(201, "创建失败")
		}
	}

	return c.Success(nil)
}

var AdminCateDelete core.HandlerFunc = func(c *core.Context) core.Response {
	id := helpers.StrTo(c.Query("id")).MustInt()
	total, _ := gosql.Model(&models.Posts{}).Where("cate_id = ?", id).Count()

	if total > 0 {
		return c.Fail(201, "该分类下面还有文章，不能删除")
	}

	if _, err := gosql.Model(&models.Cates{Id: id}).Delete(); err != nil {
		logger.Error(err)
		return c.Fail(201, "删除失败")
	}
	return c.Redirect("/admin/cates")
}
