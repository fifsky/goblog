package controllers

import (
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/models"
	"github.com/fifsky/goblog/helpers"
	"github.com/ilibs/sessions"
	"github.com/fifsky/goblog/helpers/pagination"
	"github.com/gin-gonic/gin/binding"
	"github.com/ilibs/gosql"
	"github.com/fifsky/goblog/core"
	"github.com/ilibs/logger"
	"os"
	"io"
)

func AdminIndex(c *gin.Context) {
	h := defaultH(c)
	h["Options"] = c.GetStringMapString("options")
	c.HTML(http.StatusOK, "admin/index", h)
}

func AdminIndexPost(c *gin.Context) {
	c.Request.ParseForm()
	options := c.Request.PostForm

	for k, v := range options {
		gosql.Model(&models.Options{
			OptionValue: v[0],
		}).Where("option_key = ?", k).Update()
	}

	o, err := models.GetOptions()

	if err == nil {
		core.Global.Store("options", o)
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"message":    "保存成功",
		"options":    options,
	})
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

	user := &models.Users{Name: user_name}
	err := gosql.Model(user).Get()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "用户名或者密码错误",
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

func AdminMoodGet(c *gin.Context) {
	h := defaultH(c)

	id := helpers.StrTo(c.Query("id")).MustInt()
	if id > 0 {
		mood := &models.Moods{Id: id}
		gosql.Model(mood).Get()
		h["Mood"] = mood
	}

	num := 10

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	moods, err := models.MoodGetList(page, num)
	h["Moods"] = moods

	total, err := gosql.Model(&models.Moods{}).Count()
	pager := pagination.New(int(total), num, page, 3)
	h["Pager"] = pager

	if err == nil {
		c.HTML(http.StatusOK, "admin/moods", h)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func AdminMoodPost(c *gin.Context) {
	moods := &models.Moods{}
	if err := c.Bind(moods); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "参数错误:" + err.Error(),
		})
		return
	}

	if user, exists := c.Get("LoginUser"); exists {
		moods.UserId = user.(*models.Users).Id
	}

	if moods.Content == "" {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "内容不能为空",
		})
		return
	}

	if moods.Id > 0 {
		if _, err := gosql.Model(moods).Update(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "更新失败",
			})
			logger.Error(err)
			return
		}
	} else {
		if _, err := gosql.Model(moods).Create(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "发表失败",
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

func AdminMoodDelete(c *gin.Context) {
	id := helpers.StrTo(c.Query("id")).MustInt()

	mood := &models.Moods{Id: id}
	if _, err := gosql.Model(mood).Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "删除失败",
		})
		logger.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/moods")
}

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
	cates, err := models.CateGetList(page, num)
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

func AdminLinkGet(c *gin.Context) {
	h := defaultH(c)

	id, _ := helpers.StrTo(c.Query("id")).Int()
	if id > 0 {
		link := &models.Links{Id: id}
		gosql.Model(link).Get()
		h["Link"] = link
	}

	num := 10

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	links, err := models.LinkGetList(page, num)
	h["Links"] = links

	total, err := gosql.Model(&models.Links{}).Count()
	pager := pagination.New(int(total), num, page, 3)
	h["Pager"] = pager

	if err == nil {
		c.HTML(http.StatusOK, "admin/links", h)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func AdminLinkPost(c *gin.Context) {
	links := &models.Links{}
	if err := c.Bind(links); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "参数错误:" + err.Error(),
		})
		return
	}

	if links.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "连接名称不能为空",
		})
		return
	}

	if links.Url == "" {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "连接地址不能为空",
		})
		return
	}

	if links.Id > 0 {
		if _, err := gosql.Model(links).Update(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "更新失败",
			})
			logger.Error(err)
			return
		}
	} else {
		if _, err := gosql.Model(links).Create(); err != nil {
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

func AdminLinkDelete(c *gin.Context) {
	id := helpers.StrTo(c.Query("id")).MustInt()

	if _, err := gosql.Model(&models.Links{Id: id}).Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "删除失败",
		})
		logger.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/links")
}

func AdminRemindGet(c *gin.Context) {
	h := defaultH(c)

	id := helpers.StrTo(c.Query("id")).MustInt()
	if id > 0 {
		remind := &models.Reminds{Id: id}
		gosql.Model(remind).Get()
		h["Remind"] = remind
	}

	num := 10

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	reminds, err := models.RemindGetList(page, num)

	h["Reminds"] = reminds

	total, err := gosql.Model(&models.Reminds{}).Count()
	pager := pagination.New(int(total), num, page, 3)
	h["Pager"] = pager

	h["Types"] = map[int]string{
		0: "固定",
		1: "每分钟",
		2: "每小时",
		3: "每天",
		4: "每周",
		5: "每月",
		6: "每年",
	}

	h["Layouts"] = map[int]string{
		0: "2006-01-02 15:04:05",
		1: "",
		2: "",
		3: "15:04:00",
		4: "15:04:00",
		5: "02日15:04:05",
		6: "01月02日15:04:05",
	}

	h["CurrDate"] = time.Now().Format("2006-01-02 15:04:05")

	if err == nil {
		c.HTML(http.StatusOK, "admin/remind", h)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func AdminRemindPost(c *gin.Context) {
	reminds := &models.Reminds{}
	if err := c.ShouldBindWith(reminds, binding.Form); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "参数错误:" + err.Error(),
		})
		return
	}

	if reminds.Content == "" {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "提醒内容不能为空",
		})
		return
	}

	if reminds.Id > 0 {
		if _, err := gosql.Model(reminds).Update(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "更新失败:" + err.Error(),
			})
			logger.Error(err)
			return
		}
	} else {
		if _, err := gosql.Model(reminds).Create(); err != nil {
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

func AdminRemindDelete(c *gin.Context) {
	id := helpers.StrTo(c.Query("id")).MustInt()

	if _, err := gosql.Model(&models.Reminds{Id: id}).Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "删除失败",
		})
		logger.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/remind")
}

func AdminUsersGet(c *gin.Context) {
	h := defaultH(c)
	num := 10

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	users, err := models.UserGetList(page, num)

	h["Users"] = users

	total, err := gosql.Model(&models.Users{}).Count()
	pager := pagination.New(int(total), num, page, 3)
	h["Pager"] = pager

	if err == nil {
		c.HTML(http.StatusOK, "admin/users", h)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func AdminUserGet(c *gin.Context) {
	h := defaultH(c)
	id := helpers.StrTo(c.Query("id")).MustInt()
	if id > 0 {
		user := &models.Users{Id: id}
		err := gosql.Model(user).Get()
		if err != nil {
			HandleMessage(c, "用户不存在", "您访问的用户不存在或已经删除！")
			return
		}
		h["User"] = user
	}

	c.HTML(http.StatusOK, "admin/post_user", h)
}

func AdminUserPost(c *gin.Context) {
	users := &models.Users{}
	if err := c.Bind(users); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "参数错误:" + err.Error(),
		})
		return
	}

	if users.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "用户名不能为空",
		})
		return
	}

	if users.Password == "" {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "密码不能为空",
		})
		return
	}

	if users.NickName == "" {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "昵称不能为空",
		})
		return
	}

	users.Password = helpers.Md5(users.Password)

	if users.Id > 0 {
		if _, err := gosql.Model(users).Update(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "更新失败",
			})
			logger.Error(err)
			return
		}
	} else {
		if _, err := gosql.Model(users).Create(); err != nil {
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

func AdminUserStatus(c *gin.Context) {
	id := helpers.StrTo(c.Query("id")).MustInt()
	status := helpers.StrTo(c.Query("status")).MustInt()

	if _, err := gosql.Model(&models.Users{}).Where("id = ? and status = ?", id, status).Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "删除失败",
		})
		logger.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/users")
}
