package middleware

import (
	"github.com/fifsky/goblog/core"
	"github.com/fifsky/goblog/models"
	"github.com/verystar/logger"
)

var AuthLogin core.HandlerFunc = func(c *core.Context) core.Response {
	if user, ok := c.SharedData["LoginUser"]; ok {
		if _, ok := user.(*models.Users); ok {
			c.Next()
			return nil
		}
	}

	logger.Error("User not authorized to visit %s", c.Request.RequestURI)
	return c.Redirect("/admin/login")
}
