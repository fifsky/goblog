package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/models"
	"net/http"
	"github.com/fifsky/goblog/helpers"
	"github.com/fifsky/goblog/helpers/pagination"
)

func IndexGet(c *gin.Context) {
	num := 10
	domain := c.Param("domain")
	cate := &models.Cates{}

	if domain != "" {
		cate.Domain = domain
		cate.Get()
	}

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	postModel := new(models.Posts)
	if cate.Id > 0 {
		postModel.CateId = cate.Id
	}

	posts, err := postModel.GetList(page, num)

	h := defaultH(c)
	h["Posts"] = posts

	total, err := postModel.Count()
	pager := pagination.New(int(total), num, page, 3)
	h["Pager"] = pager

	if err == nil {
		c.HTML(http.StatusOK, "index/index", h)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func ArticleGet(c *gin.Context) {
	id, _ := helpers.StrTo(c.Param("id")).Uint()
	postModel := &models.Posts{Id: id}
	post, err := postModel.Get()

	cateModel := &models.Cates{Id: post.CateId}
	cate, err := cateModel.Get()

	userModel := &models.Users{Id: post.UserId}
	user, err := userModel.Get()

	newpost := &models.UserPosts{Posts: *post, Name: cate.Name, Domain: cate.Domain, NickName: user.NickName}

	h := defaultH(c)
	h["Title"] = post.Title
	h["Post"] = newpost

	if err == nil {
		c.HTML(http.StatusOK, "index/article", h)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
