package handler

import (
	"image/png"
	"net/http"
	"os"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/fifsky/goblog/config"
	"github.com/fifsky/goblog/core"
	"github.com/fifsky/goblog/helpers"
	"github.com/fifsky/goblog/models"
	"github.com/gin-gonic/gin"
	"github.com/ilibs/gosql"
	"github.com/nfnt/resize"
	"github.com/verystar/logger"
)

var AdminArticlesGet core.HandlerFunc = func(c *core.Context) core.Response {
	num := 10
	h := gin.H{}
	post := &models.Posts{}
	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	cateId := helpers.StrTo(c.DefaultQuery("cate_id", "0")).MustInt()
	h["CateId"] = cateId
	if cateId > 0 {
		post.CateId = cateId
	}

	posts, err := models.PostGetList(post, page, num, "", "")

	if err != nil {
		return c.ErrorMessage(err)
	}

	cates := make([]*models.Cates, 0)
	gosql.Model(&cates).All()

	h["Posts"] = posts
	h["Cates"] = cates

	total, err := gosql.Model(&models.Posts{}).Count()
	h["Pager"] = c.Pagination(total, num, page)

	if err != nil {
		return c.ErrorMessage(err)
	}

	return c.HTML("admin/articles", h)
}

var AdminArticleGet core.HandlerFunc = func(c *core.Context) core.Response {
	h := gin.H{}
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

	if err != nil {
		return c.ErrorMessage(err)
	}
	return c.HTML("admin/post_article", h)
}

var AdminArticlePost core.HandlerFunc = func(c *core.Context) core.Response {
	post := &models.Posts{}
	if err := c.ShouldBind(post); err != nil {
		return c.Fail(201, "参数错误:"+err.Error())
	}

	post.UserId = c.Session().Get("UserId").(int)

	if post.Title == "" {
		return c.Fail(201, "文章标题不能为空")
	}

	if post.Id > 0 {
		if _, err := gosql.Model(post).Update(); err != nil {
			logger.Error(err)
			return c.Fail(201, "更新文章失败")
		}
	} else {
		if _, err := gosql.Model(post).Create(); err != nil {
			logger.Error(err)
			return c.Fail(201, "发表文章失败")
		}
	}

	return c.Success(post)
}

var AdminArticleDelete core.HandlerFunc = func(c *core.Context) core.Response {
	id := helpers.StrTo(c.Query("id")).MustInt()

	post := &models.Posts{Id: id}
	if _, err := gosql.Model(post).Delete(); err != nil {
		logger.Error(err)
		return c.Fail(201, "删除失败")
	}
	return c.Redirect(c.Request.Referer())
}

var AdminUploadPost core.HandlerFunc = func(c *core.Context) core.Response {
	//开发环境上传到local
	if config.App.Common.Env == "local" {
		return AdminUploadPostLocal(c)
	}

	file, _, err := c.Request.FormFile("uploadFile")
	if err != nil {
		c.Status(http.StatusBadRequest)
		return c.String("Bad request")
	}
	client, err := oss.New(config.App.OSS.Endpoint, config.App.OSS.AccessKey, config.App.OSS.AccessSecret)

	if err != nil {
		return c.JSON(gin.H{
			"jsonrpc": "2.0",
			"error": gin.H{
				"code":    101,
				"message": err.Error(),
			},
			"id": "id",
		})
	}

	bucket, _ := client.Bucket(config.App.OSS.Bucket)
	day := time.Now().Format("20060102")

	filename := "upload/"+day+"/"+helpers.Md5File(file)+".png"
	file.Seek(0,0)

	err = bucket.PutObject(filename, file)
	if err != nil {
		return c.JSON(gin.H{
			"jsonrpc": "2.0",
			"error": gin.H{
				"code":    100,
				"message": err.Error(),
			},
			"id": "id",
		})
	}

	return c.String("https://static.fifsky.com/" + filename+"!blog")
}

var AdminUploadPostLocal core.HandlerFunc = func(c *core.Context) core.Response {
	file, header, err := c.Request.FormFile("uploadFile")
	if err != nil {
		c.Status(http.StatusBadRequest)
		return c.String("Bad request")
	}
	filename := header.Filename
	day := time.Now().Format("20060102")
	dir := "static/upload/" + day
	exists, _ := helpers.PathExists(dir)
	if !exists {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			logger.Fatal(err)
			return c.JSON(gin.H{
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

		return c.JSON(gin.H{
			"jsonrpc": "2.0",
			"error": gin.H{
				"code":    100,
				"message": "Failed to create file.",
			},
			"id": "id",
		})
	}
	defer out.Close()

	img, err := png.Decode(file)
	if err != nil {
		logger.Fatal(err)
	}
	file.Close()

	m := resize.Resize(800, 0, img, resize.Lanczos3)
	err = png.Encode(out, m)

	if err != nil {
		logger.Fatal(err)
		return c.JSON(gin.H{
			"jsonrpc": "2.0",
			"error": gin.H{
				"code":    100,
				"message": "Failed to save directory.",
			},
			"id": "id",
		})
	}
	return c.String("/static/upload/" + day + "/" + filename)
}
