package router

import (
	"github.com/fifsky/goblog/core"
	"github.com/fifsky/goblog/debug"
	"github.com/fifsky/goblog/handler"
	"github.com/fifsky/goblog/router/middleware"
	"github.com/gin-gonic/gin"
)

func Route(router *gin.Engine) {
	core.SetTemplate(router)

	//中间件
	router.Use(middleware.Sessions())
	router.Use(middleware.Ginrus())
	router.Use(core.Middleware(middleware.SharedData))

	//静态文件
	router.Static("/static", "./static")

	router.NoRoute(core.Handle(handler.Handle404))
	router.GET("/", core.Handle(handler.IndexGet))
	router.GET("/about", core.Handle(handler.AboutGet))
	router.GET("/article/:id", core.Handle(handler.ArticleGet))
	router.POST("/post/comment", core.Handle(handler.CommentPost))
	router.GET("/categroy/:domain", core.Handle(handler.IndexGet))
	router.GET("/date/:year/:month", core.Handle(handler.IndexGet))
	router.GET("/search", core.Handle(handler.IndexGet))
	router.GET("/avatar", core.Handle(handler.Avatar))

	//管理后台
	admin := router.Group("/admin")
	admin.GET("/login", core.Handle(handler.LoginGet))
	admin.POST("/login", core.Handle(handler.LoginPost))
	admin.GET("/logout", core.Handle(handler.LogoutGet))

	admin.Use(middleware.AuthLogin())
	{
		//网站设置
		admin.GET("/index", core.Handle(handler.AdminIndex))
		admin.POST("/index", core.Handle(handler.AdminIndexPost))

		//文章管理
		admin.GET("/articles", core.Handle(handler.AdminArticlesGet))
		admin.GET("/post/article", core.Handle(handler.AdminArticleGet))
		admin.POST("/post/article", core.Handle(handler.AdminArticlePost))
		admin.GET("/post/article_delete", core.Handle(handler.AdminArticleDelete))
		admin.POST("/post/upload", core.Handle(handler.AdminUploadPost))

		//心情
		admin.GET("/moods", core.Handle(handler.AdminMoodGet))
		admin.POST("/moods", core.Handle(handler.AdminMoodPost))
		admin.GET("/mood_delete", core.Handle(handler.AdminMoodDelete))

		//分类
		admin.GET("/cates", core.Handle(handler.AdminCateGet))
		admin.POST("/cates", core.Handle(handler.AdminCatePost))
		admin.GET("/cate_delete", core.Handle(handler.AdminCateDelete))

		//链接
		admin.GET("/links", core.Handle(handler.AdminLinkGet))
		admin.POST("/links", core.Handle(handler.AdminLinkPost))
		admin.GET("/link_delete", core.Handle(handler.AdminLinkDelete))

		//提醒
		admin.GET("/remind", core.Handle(handler.AdminRemindGet))
		admin.POST("/remind", core.Handle(handler.AdminRemindPost))
		admin.GET("/remind_delete", core.Handle(handler.AdminRemindDelete))

		//用户
		admin.GET("/users", core.Handle(handler.AdminUsersGet))
		admin.GET("/post/user", core.Handle(handler.AdminUserGet))
		admin.POST("/post/user", core.Handle(handler.AdminUserPost))
		admin.GET("/user_status", core.Handle(handler.AdminUserStatus))
	}

	//debug handler
	debug.Route(router)
}
