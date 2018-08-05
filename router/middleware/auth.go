package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/models"
	"github.com/ilibs/logger"
)

func AuthLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get("LoginUser"); user != nil {
			if _, ok := user.(*models.Users); ok {
				c.Next()
				return
			}
		}

		logger.Error("User not authorized to visit %s", c.Request.RequestURI)

		c.Redirect(http.StatusFound, "/admin/login")
	}
}
