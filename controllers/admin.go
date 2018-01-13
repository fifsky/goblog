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
	h["Options"] = c.GetStringMapString("options")
	c.HTML(http.StatusOK, "admin/index", h)
}

func AdminIndexPost(c *gin.Context) {
	c.Request.ParseForm()
	options := c.Request.PostForm

	for k, v := range options {
		optionModel := &models.Options{}
		optionModel.OptionKey = k
		optionModel.OptionValue = v[0]
		optionModel.Update()
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

	userModel := &models.Users{Name: user_name}
	user, has := userModel.Get()

	if !has {
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
	postModel := new(models.Posts)
	posts, err := postModel.GetList(page, num, "")

	cateModel := new(models.Cates)
	cates, err := cateModel.All()

	h := defaultH(c)
	h["Posts"] = posts
	h["Cates"] = cates

	total, err := postModel.Count("")
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

func AdminMoodGet(c *gin.Context) {
	h := defaultH(c)

	id, _ := helpers.StrTo(c.Query("id")).Uint()
	if id > 0 {
		moodModel := &models.Moods{Id: id}
		mood, _ := moodModel.Get()
		h["Mood"] = mood
	}

	num := 10

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	model := new(models.Moods)
	moods, err := model.GetList(page, num)

	h["Moods"] = moods

	total, err := model.Count()
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
		if _, err := moods.Update(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "更新失败",
			})
			logrus.Error(err)
			return
		}
	} else {
		if _, err := moods.Insert(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "发表失败",
			})
			logrus.Error(err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"message":    "ok",
	})
}

func AdminMoodDelete(c *gin.Context) {
	id, _ := helpers.StrTo(c.Query("id")).Uint()

	mood := &models.Moods{Id: id}
	if _, err := mood.Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "删除失败",
		})
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/moods")
}

func AdminCateGet(c *gin.Context) {
	h := defaultH(c)

	id, _ := helpers.StrTo(c.Query("id")).Uint()
	if id > 0 {
		cateModel := &models.Cates{Id: id}
		cate, _ := cateModel.Get()
		h["Cate"] = cate
	}

	num := 10

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	model := new(models.Cates)
	cates, err := model.GetList(page, num)

	h["Cates"] = cates

	total, err := model.Count()
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
		if _, err := cates.Update(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "更新失败",
			})
			logrus.Error(err)
			return
		}
	} else {
		if _, err := cates.Insert(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "创建失败",
			})
			logrus.Error(err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"message":    "ok",
	})
}

func AdminCateDelete(c *gin.Context) {
	id, _ := helpers.StrTo(c.Query("id")).Uint()

	post := &models.Posts{CateId: id}
	total, _ := post.Count("")

	if total > 0 {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "该分类下面还有文章，不能删除",
		})
		return
	}

	mood := &models.Cates{Id: id}
	if _, err := mood.Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "删除失败",
		})
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/cates")
}

func AdminLinkGet(c *gin.Context) {
	h := defaultH(c)

	id, _ := helpers.StrTo(c.Query("id")).Uint()
	if id > 0 {
		linkModel := &models.Links{Id: id}
		link, _ := linkModel.Get()
		h["Link"] = link
	}

	num := 10

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	model := new(models.Links)
	links, err := model.GetList(page, num)

	h["Links"] = links

	total, err := model.Count()
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
		if _, err := links.Update(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "更新失败",
			})
			logrus.Error(err)
			return
		}
	} else {
		if _, err := links.Insert(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "创建失败",
			})
			logrus.Error(err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"message":    "ok",
	})
}

func AdminLinkDelete(c *gin.Context) {
	id, _ := helpers.StrTo(c.Query("id")).Uint()

	link := &models.Links{Id: id}
	if _, err := link.Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "删除失败",
		})
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/links")
}


func AdminRemindGet(c *gin.Context) {
	h := defaultH(c)

	id, _ := helpers.StrTo(c.Query("id")).Uint()
	if id > 0 {
		linkModel := &models.Links{Id: id}
		link, _ := linkModel.Get()
		h["Link"] = link
	}

	num := 10

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	model := new(models.Links)
	links, err := model.GetList(page, num)

	h["Links"] = links

	total, err := model.Count()
	pager := pagination.New(int(total), num, page, 3)
	h["Pager"] = pager

	if err == nil {
		c.HTML(http.StatusOK, "admin/links", h)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func AdminRemindPost(c *gin.Context) {
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
		if _, err := links.Update(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "更新失败",
			})
			logrus.Error(err)
			return
		}
	} else {
		if _, err := links.Insert(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "创建失败",
			})
			logrus.Error(err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"message":    "ok",
	})
}

func AdminRemindDelete(c *gin.Context) {
	id, _ := helpers.StrTo(c.Query("id")).Uint()

	link := &models.Links{Id: id}
	if _, err := link.Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "删除失败",
		})
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/links")
}


func AdminUsersGet(c *gin.Context) {
	h := defaultH(c)
	num := 10

	page := helpers.StrTo(c.DefaultQuery("page", "1")).MustInt()
	model := new(models.Users)
	users, err := model.GetList(page, num)

	h["Users"] = users

	total, err := model.Count()
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
	id, _ := helpers.StrTo(c.Query("id")).Uint()
	if id > 0 {
		userModel := &models.Users{Id: id}
		user, has := userModel.Get()
		if !has {
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
		if _, err := users.Update(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "更新失败",
			})
			logrus.Error(err)
			return
		}
	} else {
		if _, err := users.Insert(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 201,
				"message":    "创建失败",
			})
			logrus.Error(err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"message":    "ok",
	})
}

func AdminUserStatus(c *gin.Context) {
	id, _ := helpers.StrTo(c.Query("id")).Uint()
	status, _ := helpers.StrTo(c.Query("status")).Uint8()

	user := &models.Users{Id: id, Status: status}
	if _, err := user.Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 201,
			"message":    "删除失败",
		})
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/users")
}
