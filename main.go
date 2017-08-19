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

	// log to file
	f, ferr := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if ferr != nil {
		fmt.Printf("error opening file: %v", ferr)
	}
	defer f.Close()

	logrus.SetFormatter(&logrus.TextFormatter{})
	//logrus.SetOutput(os.Stdout)
	logrus.SetOutput(f)
	logrus.SetLevel(logrus.WarnLevel)

	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.Use(ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true))

	system.LoadConfig()
	_, err := models.InitDB()
	if err != nil {
		logrus.Error("err open databases", err)
		return
	}

	setTemplate(router)
	setSessions(router)

	router.Use(SharedData())

	router.Static("/static", "./static")

	router.NoRoute(controllers.Handle404)
	router.GET("/", controllers.IndexGet)

	visitor := router.Group("/visitor")
	visitor.Use(AuthRequired())
	{
	}

	router.GET("/article/:id", controllers.ArchiveGet)

	router.Run(":8080")
}

func setTemplate(engine *gin.Engine) {

	funcMap := template.FuncMap{
		"DateFormat": helpers.DateFormat,
		"Substr":     helpers.Substr,
		"Truncate":   helpers.Truncate,
		"Unescaped":  helpers.Unescaped,
		"StaticUrl":  helpers.StaticUrl,
	}

	engine.SetFuncMap(funcMap)
	engine.LoadHTMLGlob("views/**/*")
}

//setSessions initializes sessions & csrf middlewares
func setSessions(router *gin.Engine) {
	config := system.GetConfig()
	//https://github.com/gin-gonic/contrib/tree/master/sessions
	store := sessions.NewCookieStore([]byte(config.SessionSecret))
	store.Options(sessions.Options{HttpOnly: true, MaxAge: 7 * 86400, Path: "/"}) //Also set Secure: true if using SSL, you should though
	router.Use(sessions.Sessions("gin-session", store))
}

//+++++++++++++ middlewares +++++++++++++++++++++++

//SharedData fills in common data, such as user info, etc...
func SharedData() gin.HandlerFunc {
	return func(c *gin.Context) {

		optionModel := &models.Options{}
		options, _ := optionModel.GetOptions()

		c.Set("options", options)

		session := sessions.Default(c)
		if id := session.Get(controllers.SESSION_KEY); id != nil {
			//user := models.DB.First(&models.User{}, id)
			//if user.Error == nil {
			//	c.Set(controllers.CONTEXT_USER_KEY, user)
			//}
		}
		c.Next()
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get(controllers.CONTEXT_USER_KEY); user != nil {
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
