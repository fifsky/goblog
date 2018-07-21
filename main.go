package main

import (
	"net/http"
	"os"
	"time"
	"fmt"
	"flag"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/fifsky/goblog/core"
	"github.com/fifsky/goblog/core/config"
	"github.com/fifsky/goblog/controllers"
	"github.com/fifsky/goblog/models"
	"github.com/fifsky/goblog/helpers"
	"github.com/ilibs/gosql"
)

func main() {
	config.LoadConfig()
	fmt.Println("Run Mode:", gin.Mode())
	gosql.Connect(config.GetConfig().Database)

	flag.Parse()
	cmd := flag.Arg(0)
	if cmd == "install" {
		_, err := config.ImportDB()
		if err != nil {
			fmt.Println("Import DB Error:" + err.Error())
			logrus.Error(err)
		}
		return
	}

	f := setLogger()
	defer f.Close()

	router := gin.Default()
	router.Use(core.Ginrus(logrus.StandardLogger(), time.RFC3339, true))
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

	admin.Use(authLogin())
	{
		//网站设置
		admin.GET("/index", controllers.AdminIndex)
		admin.POST("/index", controllers.AdminIndexPost)

		//文章管理
		admin.GET("/articles", controllers.AdminArticlesGet)
		admin.GET("/post/article", controllers.AdminArticleGet)
		admin.POST("/post/article", controllers.AdminArticlePost)
		admin.GET("/post/article_delete", controllers.AdminArticleDelete)

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

func setLogger() *os.File {
	f, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	logrus.SetFormatter(&logrus.TextFormatter{})
	//logrus.SetOutput(os.Stdout)
	logrus.SetOutput(f)
	if gin.Mode() == gin.DebugMode {
		logrus.SetLevel(logrus.InfoLevel)
	}
	return f
}

func authLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get("LoginUser"); user != nil {
			if _, ok := user.(*models.Users); ok {
				c.Next()
				return
			}
		}

		logrus.Warnf("User not authorized to visit %s", c.Request.RequestURI)

		c.Redirect(http.StatusFound, "/admin/login")
	}
}
