package handler

import (
	"net/http"

	"github.com/fifsky/goblog/models"
	"github.com/fifsky/goblog/helpers"
	"github.com/ilibs/gosql"
	"github.com/fifsky/goblog/core"
	"github.com/ilibs/identicon"
)

var IndexGet core.HandlerFunc = func(c *core.Context) core.Response {
	options := c.GetStringMapString("options")
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

	keyword := c.Query("keyword")

	post.Type = 1
	posts, err := models.PostGetList(post, page, num, artdate, keyword)
	h := defaultH(c.Context)
	h["Posts"] = posts

	builder := gosql.Model(post)

	if artdate != "" {
		builder.Where("DATE_FORMAT(created_at,'%Y-%m') = ?", artdate)
	}

	if keyword != "" {
		builder.Where("title like ?", "%"+keyword+"%")
	}

	total, err := builder.Count()
	h["Pager"] = c.Pagination(total, num, page)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return nil
	}
	return c.HTML("index/index", h)
}

var ArticleGet core.HandlerFunc = func(c *core.Context) core.Response {
	id := helpers.StrTo(c.Param("id")).MustInt()
	url := c.GetString("url")
	post := &models.Posts{Id: id}

	if url != "" {
		post.Url = url
	}

	err := gosql.Model(post).Get()

	if err != nil {
		return HandleMessage(c, "文章不存在", "您访问的文章不存在或已经删除！")
	}

	cate := &models.Cates{Id: post.CateId}
	gosql.Model(cate).Get()

	user := &models.Users{Id: post.UserId}
	gosql.Model(user).Get()

	newpost := &models.UserPosts{Posts: *post, Name: cate.Name, Domain: cate.Domain, NickName: user.NickName}

	h := defaultH(c.Context)
	h["Title"] = post.Title
	h["Post"] = newpost

	if url == "" {
		prev, err := models.PostPrev(post.Id)
		if err == nil {
			h["Prev"] = prev
		}
		next, err := models.PostNext(post.Id)
		if err == nil {
			h["Next"] = next
		}
	}

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	comments, err := models.PostComments(post.Id, page, 100)
	if err == nil {
		h["Comments"] = comments
	}

	return c.HTML("index/article", h)
}

var AboutGet core.HandlerFunc = func(c *core.Context) core.Response {
	c.Set("url", "about")
	return ArticleGet(c)
}

var Avatar core.HandlerFunc = func(c *core.Context) core.Response {
	name := c.DefaultQuery("name", "default")

	// New Generator: Rehuse
	ig, err := identicon.New(
		"fifsky", // Namespace
		5,        // Number of blocks (Size)
		5,        // Density
	)

	if err != nil {
		panic(err) // Invalid Size or Density
	}

	ii, err := ig.Draw(name) // Generate an IdentIcon

	if err != nil {
		return nil
	}
	// Takes the size in pixels and any io.Writer
	ii.Png(300, c.Writer) // 300px * 300px
	return nil
}
