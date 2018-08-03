package route

import (
	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/core"
	"github.com/fifsky/goblog/handler"
	"github.com/fifsky/goblog/debug"
	"github.com/fifsky/goblog/route/middleware"
)

func Route(router *gin.Engine)  {
	core.SetTemplate(router)

	//中间件
	router.Use(middleware.Sessions())
	router.Use(middleware.Ginrus())
	router.Use(middleware.SharedData())

	//静态文件
	router.Static("/static", "./static")

	router.NoRoute(handler.Handle404)
	router.GET("/", handler.IndexGet)
	router.GET("/about", handler.AboutGet)
	router.GET("/article/:id", handler.ArticleGet)
	router.GET("/categroy/:domain", handler.IndexGet)
	router.GET("/date/:year/:month", handler.IndexGet)
	router.GET("/search", handler.IndexGet)


	//管理后台
	admin := router.Group("/admin")
	admin.GET("/login", core.Handle(handler.LoginGet))
	admin.POST("/login", handler.LoginPost)
	admin.GET("/logout", handler.LogoutGet)

	admin.Use(middleware.AuthLogin())
	{
		//网站设置
		admin.GET("/index", handler.AdminIndex)
		admin.POST("/index", handler.AdminIndexPost)

		//文章管理
		admin.GET("/articles", handler.AdminArticlesGet)
		admin.GET("/post/article", handler.AdminArticleGet)
		admin.POST("/post/article", handler.AdminArticlePost)
		admin.GET("/post/article_delete", handler.AdminArticleDelete)
		admin.POST("/post/upload", handler.AdminUploadPost)


		//心情
		admin.GET("/moods", handler.AdminMoodGet)
		admin.POST("/moods", handler.AdminMoodPost)
		admin.GET("/mood_delete", handler.AdminMoodDelete)

		//分类
		admin.GET("/cates", handler.AdminCateGet)
		admin.POST("/cates", handler.AdminCatePost)
		admin.GET("/cate_delete", handler.AdminCateDelete)

		//链接
		admin.GET("/links", handler.AdminLinkGet)
		admin.POST("/links", handler.AdminLinkPost)
		admin.GET("/link_delete", handler.AdminLinkDelete)

		//提醒
		admin.GET("/remind", handler.AdminRemindGet)
		admin.POST("/remind", handler.AdminRemindPost)
		admin.GET("/remind_delete", handler.AdminRemindDelete)

		//用户
		admin.GET("/users", handler.AdminUsersGet)
		admin.GET("/post/user", handler.AdminUserGet)
		admin.POST("/post/user", handler.AdminUserPost)
		admin.GET("/user_status", handler.AdminUserStatus)
	}

	debugger := router.Group("/debug")
	{
		debugger.Use(middleware.AuthLogin())
		debugger.GET("/info", debug.InfoHandler())
		debugger.GET("/pprof/", debug.IndexHandler())
		debugger.GET("/pprof/heap", debug.HeapHandler())
		debugger.GET("/pprof/goroutine", debug.GoroutineHandler())
		debugger.GET("/pprof/block", debug.BlockHandler())
		debugger.GET("/pprof/threadcreate", debug.ThreadCreateHandler())
		debugger.GET("/pprof/cmdline", debug.CmdlineHandler())
		debugger.GET("/pprof/profile", debug.ProfileHandler())
		debugger.GET("/pprof/symbol", debug.SymbolHandler())
		debugger.POST("/pprof/symbol", debug.SymbolHandler())
		debugger.GET("/pprof/trace", debug.TraceHandler())
	}
}
