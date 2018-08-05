package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/models"
	"github.com/ilibs/gosql"
	"github.com/ilibs/sessions"
	"github.com/fifsky/goblog/core"
)

//middlewares
func SharedData() gin.HandlerFunc {
	return func(c *gin.Context) {

		if !strings.HasPrefix(c.Request.URL.Path, "/static") {
			//网站全局配置
			options, ok := core.Global.Load("options")
			if !ok {
				options, _ = models.GetOptions()
				core.Global.Store("options", options)
				c.Set("options", options)
			} else {
				c.Set("options", options.(map[string]string))
			}

			session := sessions.Default(c)
			if uid := session.Get("UserId"); uid != nil {
				if user, ok := core.Global.Load("LoginUser"); ok {
					c.Set("LoginUser", user.(*models.Users))
				} else {
					user = &models.Users{}
					err := gosql.Model(user).Where("id = ?", uid).Get()
					if err == nil {
						core.Global.Store("LoginUser", user)
						c.Set("LoginUser", user)
					}
				}
			}
		}

		c.Next()
	}
}
