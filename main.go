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
)

func main() {
	gin.SetMode(gin.DebugMode)
	f := setLogger()
	defer f.Close()
	system.LoadConfig()
	connectDB()

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
	router.GET("/article/:id", controllers.ArchiveGet)

	admin := router.Group("/admin")
	admin.GET("/login", controllers.LoginGet)
	admin.POST("/login", controllers.LoginPost)
	admin.GET("/logout", controllers.LogoutGet)

	admin.Use(authLogin())
	{
		admin.GET("/index", controllers.AdminIndex)
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
	return func(c *gin.Context) {

		//网站全局配置
		optionModel := &models.Options{}
		options, _ := optionModel.GetOptions()
		c.Set("options", options)

		session := sessions.Default(c)
		if uID := session.Get("UserId"); uID != nil {
			userModel := &models.Users{Id: uID.(uint)}
			user, _ := userModel.Get()
			if user.Id != 0 {
				c.Set("LoginUser", user)
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
