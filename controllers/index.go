package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/models"
	"net/http"
	"github.com/fifsky/goblog/helpers"
	"github.com/fifsky/goblog/helpers/pagination"
)

func IndexGet(c *gin.Context) {
	options := c.MustGet("options").(map[string]string)
	num, err := helpers.StrTo(options["post_num"]).Int()
	if err != nil || num < 1 {
		num = 10
	}

	domain := c.Param("domain")
	cate := &models.Cates{}

	if domain != "" {
		cate.Domain = domain
		cate.Get()
	}

	artdate := ""
	year := c.Param("year")
	month := c.Param("month")

	if year != "" && month != "" {
		artdate = year + "-" + month
	}

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	postModel := new(models.Posts)
	if cate.Id > 0 {
		postModel.CateId = cate.Id
	}

	postModel.Type = 1
	posts, err := postModel.GetList(page, num, artdate)

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
	url := c.GetString("url")
	postModel := &models.Posts{Id: id}

	if url != "" {
		postModel.Url = url
	}

	post, has := postModel.Get()

	if !has {
		HandleMessage(c, "文章不存在", "您访问的文章不存在或已经删除！")
		return
	}

	cateModel := &models.Cates{Id: post.CateId}
	cate, _ := cateModel.Get()

	userModel := &models.Users{Id: post.UserId}
	user, _ := userModel.Get()

	newpost := &models.UserPosts{Posts: *post, Name: cate.Name, Domain: cate.Domain, NickName: user.NickName}

	h := defaultH(c)
	h["Title"] = post.Title
	h["Post"] = newpost

	if url == "" {
		prev, has := postModel.Prev(post.Id)
		if has {
			h["Prev"] = prev
		}
		next, has := postModel.Next(post.Id)
		if has {
			h["Next"] = next
		}
	}

	c.HTML(http.StatusOK, "index/article", h)
}

func AboutGet(c *gin.Context) {
	c.Set("url", "about")
	ArticleGet(c)
}
