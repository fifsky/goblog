package main

import (
	"github.com/sirupsen/logrus"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/controllers"
	"github.com/fifsky/goblog/helpers"
	"github.com/fifsky/goblog/models"
	"github.com/fifsky/goblog/system"
	"html/template"
	"net/http"
	"os"
	"time"
	"fmt"
	"flag"
	"strings"
)

func main() {
	system.LoadConfig()
	connectDB()

	flag.Parse()
	cmd := flag.Arg(0)
	if cmd == "install" {
		_, err := models.ImportDB()
		if err != nil {
			fmt.Println("Import DB Error:" + err.Error())
			logrus.Error(err)
		}
		return
	}

	gin.SetMode(gin.DebugMode)
	f := setLogger()
	defer f.Close()

	router := gin.Default()
	router.Use(ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true))
	setTemplate(router)
	setSessions(router)

	//中间件
	router.Use(sharedData())

	//静态文件
	router.Static("/static", "./static")

	router.NoRoute(controllers.Handle404)
	router.GET("/", controllers.IndexGet)
	router.GET("/about", controllers.AboutGet)
	router.GET("/article/:id", controllers.ArticleGet)
	router.GET("/categroy/:domain", controllers.IndexGet)

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

		//用户
		admin.GET("/users", controllers.AdminUsersGet)
		admin.GET("/post/user", controllers.AdminUserGet)
		admin.POST("/post/user", controllers.AdminUserPost)
		admin.GET("/user_status", controllers.AdminUserStatus)
	}
	router.Run(":8080")
}

func connectDB() {
	_, err := models.InitDB()
	if err != nil {
		logrus.Error("err open databases", err)
		panic(err)
	}
}

func setTemplate(engine *gin.Engine) {

	funcMap := template.FuncMap{
		"DateFormat": helpers.DateFormat,
		"Substr":     helpers.Substr,
		"Truncate":   helpers.Truncate,
		"Unescaped":  helpers.Unescaped,
		"StaticUrl":  helpers.StaticUrl,
		"IsPage":     helpers.IsPage,
	}

	engine.SetFuncMap(funcMap)
	engine.LoadHTMLGlob("views/**/*")
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

func setSessions(router *gin.Engine) {
	config := system.GetConfig()
	store := sessions.NewCookieStore([]byte(config.SessionSecret))
	store.Options(sessions.Options{HttpOnly: true, MaxAge: 7 * 86400, Path: "/"}) //Also set Secure: true if using SSL, you should though
	router.Use(sessions.Sessions("gin-session", store))

	//https://github.com/utrack/gin-csrf
	//router.Use(csrf.Middleware(csrf.Options{
	//	Secret: config.SessionSecret,
	//	ErrorFunc: func(c *gin.Context) {
	//		c.String(400, "CSRF token mismatch")
	//		c.Abort()
	//	},
	//}))
}

//middlewares
func sharedData() gin.HandlerFunc {
	global := &gin.Context{}

	return func(c *gin.Context) {

		if !strings.HasPrefix(c.Request.URL.Path, "/static") {
			//网站全局配置
			options := global.GetStringMapString("options")
			if len(options) == 0 {
				optionModel := &models.Options{}
				options, _ = optionModel.GetOptions()
				global.Set("options", options)

			}
			c.Set("options", options)

			session := sessions.Default(c)
			if uid := session.Get("UserId"); uid != nil {
				userI, exists := global.Get("LoginUser")
				var user *models.Users
				if exists {
					user = userI.(*models.Users)
				} else {
					userModel := &models.Users{Id: uid.(uint)}
					user, _ = userModel.Get()
					global.Set("LoginUser", user)
				}

				if user.Id != 0 {
					c.Set("LoginUser", user)
				}
			}
		}
		c.Next()
	}
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
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Forbidden!",
		})
		c.Abort()
	}
}
