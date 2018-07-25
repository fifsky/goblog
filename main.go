package main

import (
	"os"
	"fmt"
	"flag"
	"io/ioutil"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/core"
	"github.com/fifsky/goblog/core/config"
	"github.com/fifsky/goblog/controllers"
	"github.com/fifsky/goblog/helpers"
	"github.com/ilibs/gosql"
	"github.com/ilibs/logger"
)

func main() {
	fmt.Println("Run Mode:", gin.Mode())
	gosql.Connect(config.App.DB)

	flag.Parse()
	cmd := flag.Arg(0)
	if cmd == "install" {
		_, err := config.ImportDB()
		if err != nil {
			fmt.Println("Import DB Error:" + err.Error())
			logger.Error(err)
		}
		return
	}

	router := gin.Default()
	router.Use(core.Ginrus())
	helpers.SetTemplate(router)
	core.SetSessions(router)

	//中间件
	router.Use(core.SharedData())

	//静态文件
	router.Static("/static", "./static")

	router.NoRoute(controllers.Handle404)
	router.GET("/", controllers.IndexGet)
	router.GET("/about", controllers.AboutGet)
	router.GET("/article/:id", controllers.ArticleGet)
	router.GET("/categroy/:domain", controllers.IndexGet)
	router.GET("/date/:year/:month", controllers.IndexGet)

	//管理后台
	admin := router.Group("/admin")
	admin.GET("/login", controllers.LoginGet)
	admin.POST("/login", controllers.LoginPost)
	admin.GET("/logout", controllers.LogoutGet)

	admin.Use(core.AuthLogin())
	{
		//网站设置
		admin.GET("/index", controllers.AdminIndex)
		admin.POST("/index", controllers.AdminIndexPost)

		//文章管理
		admin.GET("/articles", controllers.AdminArticlesGet)
		admin.GET("/post/article", controllers.AdminArticleGet)
		admin.POST("/post/article", controllers.AdminArticlePost)
		admin.GET("/post/article_delete", controllers.AdminArticleDelete)
		admin.POST("/post/upload", controllers.AdminUploadPost)


		//心情
		admin.GET("/moods", controllers.AdminMoodGet)
		admin.POST("/moods", controllers.AdminMoodPost)
		admin.GET("/mood_delete", controllers.AdminMoodDelete)

		//分类
		admin.GET("/cates", controllers.AdminCateGet)
		admin.POST("/cates", controllers.AdminCatePost)
		admin.GET("/cate_delete", controllers.AdminCateDelete)

		//链接
		admin.GET("/links", controllers.AdminLinkGet)
		admin.POST("/links", controllers.AdminLinkPost)
		admin.GET("/link_delete", controllers.AdminLinkDelete)

		//提醒
		admin.GET("/remind", controllers.AdminRemindGet)
		admin.POST("/remind", controllers.AdminRemindPost)
		admin.GET("/remind_delete", controllers.AdminRemindDelete)

		//用户
		admin.GET("/users", controllers.AdminUsersGet)
		admin.GET("/post/user", controllers.AdminUserGet)
		admin.POST("/post/user", controllers.AdminUserPost)
		admin.GET("/user_status", controllers.AdminUserStatus)
	}

	setPid(os.Getpid())

	go core.StartCron()

	router.Run(":8080")
}

func setPid(pid int) {
	d := []byte(strconv.Itoa(pid))
	err := ioutil.WriteFile("./blog.pid", d, 0644)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
}
