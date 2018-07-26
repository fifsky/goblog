package controllers

import (
	"os"
	"io"
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/models"
	"github.com/fifsky/goblog/helpers"
	"github.com/fifsky/goblog/helpers/pagination"
	"github.com/ilibs/gosql"
	"github.com/ilibs/logger"
)

func AdminArticlesGet(c *gin.Context) {

	num := 10

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	posts, err := models.PostGetList(&models.Posts{}, page, num, "")

	cates := make([]*models.Cates, 0)
	gosql.Model(&cates).All()

	h := defaultH(c)
	h["Posts"] = posts
	h["Cates"] = cates

	total, err := gosql.Model(&models.Posts{}).Count()
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
	id := helpers.StrTo(c.Query("id")).MustInt()

	if id > 0 {
		post := &models.Posts{Id: id}
		gosql.Model(post).Get()
		cate := &models.Cates{Id: post.CateId}
		gosql.Model(cate).Get()
		user := &models.Users{Id: post.UserId}
		gosql.Model(user).Get()

		newpost := &models.UserPosts{Posts: *post, Name: cate.Name, Domain: cate.Domain, NickName: user.NickName}
		h["Post"] = newpost
	}

	cates := make([]*models.Cates, 0)
	err := gosql.Model(&cates).All()
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
		if _, err := gosql.Model(post).Update(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "更新文章失败",
				"post":       post,
			})
			logger.Error(err)
			return
		}
	} else {
		if _, err := gosql.Model(post).Create(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "发表文章失败",
				"post":       post,
			})
			logger.Error(err)
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
	id := helpers.StrTo(c.Query("id")).MustInt()

	post := &models.Posts{Id: id}
	if _, err := gosql.Model(post).Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "删除失败",
			"post":       post,
		})
		logger.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/articles")
}

func AdminUploadPost(c *gin.Context) {
	file, header, err := c.Request.FormFile("wangEditorPasteFile")
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}
	filename := header.Filename
	day := time.Now().Format("20060102")
	dir := "static/upload/" + day
	exists, _ := helpers.PathExists(dir)
	if !exists {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			logger.Fatal(err)
			c.JSON(http.StatusOK, gin.H{
				"jsonrpc": "2.0",
				"error": gin.H{
					"code":    100,
					"message": "Failed to create directory.",
				},
				"id": "id",
			})
		}
	}

	out, err := os.Create(dir + "/" + filename)
	if err != nil {
		logger.Fatal(err)

		c.JSON(http.StatusOK, gin.H{
			"jsonrpc": "2.0",
			"error": gin.H{
				"code":    100,
				"message": "Failed to create file.",
			},
			"id": "id",
		})
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		logger.Fatal(err)
		c.JSON(http.StatusOK, gin.H{
			"jsonrpc": "2.0",
			"error": gin.H{
				"code":    100,
				"message": "Failed to save directory.",
			},
			"id": "id",
		})
	}
	c.String(http.StatusOK, "//"+c.Request.Host+"/static/upload/"+day+"/"+filename)
}
