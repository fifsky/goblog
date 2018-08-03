package middleware

import (
	"sync"
	"github.com/gin-gonic/gin"
	"strings"
	"github.com/fifsky/goblog/models"
	"github.com/ilibs/gosql"
	"github.com/ilibs/sessions"
)

var Global = sync.Map{}

//middlewares
func SharedData() gin.HandlerFunc {
	return func(c *gin.Context) {

		if !strings.HasPrefix(c.Request.URL.Path, "/static") {
			//网站全局配置
			options, ok := Global.Load("options")
			if !ok {
				options, _ = models.GetOptions()
				Global.Store("options", options)
				c.Set("options", options)
			} else {
				c.Set("options", options.(map[string]string))
			}

			session := sessions.Default(c)
			if uid := session.Get("UserId"); uid != nil {
				if user, ok := Global.Load("LoginUser"); ok {
					c.Set("LoginUser", user.(*models.Users))
				} else {
					user = &models.Users{}
					err := gosql.Model(user).Where("id = ?", uid).Get()
					if err == nil {
						Global.Store("LoginUser", user)
						c.Set("LoginUser", user)
					}
				}
			}
		}

		c.Next()
	}
}
