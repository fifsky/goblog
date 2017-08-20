package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/fifsky/goblog/models"
	"github.com/fifsky/goblog/helpers"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/fifsky/goblog/helpers/pagination"
	"github.com/sirupsen/logrus"
)

func AdminIndex(c *gin.Context) {
	h := defaultH(c)
	h["options"] = c.MustGet("options").(map[string]string)
	c.HTML(http.StatusOK, "admin/index", h)
}

func LoginGet(c *gin.Context) {
	h := defaultH(c)
	if h["LoginUser"] != nil {
		c.Redirect(http.StatusFound, "/admin/index")
		return
	}

	c.HTML(http.StatusOK, "admin/login", h)
}

func LogoutGet(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("UserId")
	session.Save()
	c.Redirect(http.StatusFound, "/")
}

func LoginPost(c *gin.Context) {
	session := sessions.Default(c)

	user_name := c.PostForm("user_name")
	password := c.PostForm("user_pass")

	if user_name == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "用户名或密码不能为空",
		})
		return
	}

	userModel := &models.Users{Name: user_name}
	user, err := userModel.Get()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "用户不存在:" + err.Error(),
		})
		return
	}

	if user.Password != helpers.Md5(password) {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "密码错误",
		})
		return
	}

	session.Set("UserId", user.Id)
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"message":    "ok",
	})
}

func AdminArticlesGet(c *gin.Context) {

	num := 10

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	postModel := new(models.Posts)
	posts, err := postModel.GetList(page, num)

	cateModel := new(models.Cates)
	cates, err := cateModel.All()

	h := defaultH(c)
	h["Posts"] = posts
	h["Cates"] = cates

	total, err := postModel.Count()
	pager := pagination.New(int(total), num, page, 3)
	h["Pager"] = pager

	if err == nil {
		c.HTML(http.StatusOK, "admin/articles", h)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func AdminArticleGet(c *gin.Context) {
	h := defaultH(c)
	id, _ := helpers.StrTo(c.Query("id")).Uint()

	if id > 0 {
		postModel := &models.Posts{Id: id}
		post, _ := postModel.Get()
		cateModel := &models.Cates{Id: post.CateId}
		cate, _ := cateModel.Get()

		userModel := &models.Users{Id: post.UserId}
		user, _ := userModel.Get()
		newpost := &models.UserPosts{Posts: *post, Name: cate.Name, Domain: cate.Domain, NickName: user.NickName}
		h["Post"] = newpost
	}

	cateModel2 := &models.Cates{}
	cates, err := cateModel2.All()
	h["Cates"] = cates

	if err == nil {
		c.HTML(http.StatusOK, "admin/post_article", h)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func AdminArticlePost(c *gin.Context) {
	post := &models.Posts{}
	if err := c.Bind(post); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "参数错误:" + err.Error(),
		})
		return
	}

	if user, exists := c.Get("LoginUser"); exists {
		post.UserId = user.(*models.Users).Id
	}

	if post.Title == "" {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "文章标题不能为空",
			"post":       post,
		})
		return
	}

	if post.Id > 0 {
		if _, err := post.Update(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "更新文章失败",
				"post":       post,
			})
			logrus.Error(err)
			return
		}
	} else {
		if _, err := post.Insert(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "发表文章失败",
				"post":       post,
			})
			logrus.Error(err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"message":    "ok",
		"post":       post,
	})
}

func AdminArticleDelete(c *gin.Context) {
	id, _ := helpers.StrTo(c.Query("id")).Uint()

	post := &models.Posts{Id: id}
	if _, err := post.Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "删除失败",
			"post":       post,
		})
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/articles")
}
