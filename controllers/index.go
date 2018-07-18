package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/models"
	"github.com/fifsky/goblog/helpers"
	"github.com/fifsky/goblog/helpers/pagination"
	"github.com/ilibs/gosql"
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
		gosql.Model(cate).Get()
	}

	artdate := ""
	year := c.Param("year")
	month := c.Param("month")

	if year != "" && month != "" {
		artdate = year + "-" + month
	}

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	post := &models.Posts{}
	if cate.Id > 0 {
		post.CateId = cate.Id
	}

	post.Type = 1
	posts, err := models.PostGetList(post,page, num, artdate)
	h := defaultH(c)
	h["Posts"] = posts

	builder := gosql.Model(post)

	if artdate != "" {
		builder.Where("DATE_FORMAT(created_at,'%Y-%m') = ?",artdate)
	}

	total, err := builder.Count()
	pager := pagination.New(int(total), num, page, 3)
	h["Pager"] = pager

	if err == nil {
		c.HTML(http.StatusOK, "index/index", h)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func ArticleGet(c *gin.Context) {
	id := helpers.StrTo(c.Param("id")).MustInt()
	url := c.GetString("url")
	post := &models.Posts{Id: id}

	if url != "" {
		post.Url = url
	}

	err := gosql.Model(post).Get()

	if err != nil {
		HandleMessage(c, "文章不存在", "您访问的文章不存在或已经删除！")
		return
	}

	cate := &models.Cates{Id: post.CateId}
	gosql.Model(cate).Get()

	user := &models.Users{Id: post.UserId}
	gosql.Model(user).Get()

	newpost := &models.UserPosts{Posts: *post, Name: cate.Name, Domain: cate.Domain, NickName: user.NickName}

	h := defaultH(c)
	h["Title"] = post.Title
	h["Post"] = newpost

	if url == "" {
		prev, err := models.PostPrev(post.Id)
		if err != nil {
			h["Prev"] = prev
		}
		next, err := models.PostNext(post.Id)
		if err != nil {
			h["Next"] = next
		}
	}

	c.HTML(http.StatusOK, "index/article", h)
}

func AboutGet(c *gin.Context) {
	c.Set("url", "about")
	ArticleGet(c)
}
