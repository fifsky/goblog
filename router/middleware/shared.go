package middleware

import (
	"runtime"
	"strings"

	"github.com/fifsky/goblog/core"
	"github.com/fifsky/goblog/models"
	"github.com/verystar/logger"
)

//middlewares

var SharedData core.HandlerFunc = func(c *core.Context) core.Response {
	if !strings.HasPrefix(c.Request.URL.Path, "/static") {
		if uid := c.Session().Get("UserId"); uid != nil {
			user, err := models.GetUser(uid.(int))
			if err == nil {
				c.SharedData["LoginUser"] = user
			}
		}

		//global shared data
		options, _ := models.GetOptions()
		c.Set("options", options)
		c.SharedData["Options"] = options
		c.SharedData["GOVERSION"] = runtime.Version()

		url := strings.Split(c.Request.URL.Path, "/")
		c.SharedData["UrlPath"] = c.Request.URL.Path
		c.SharedData["URI"] = c.Request.RequestURI

		if url[1] != "admin" {
			mood, err := models.MoodGetList(1,10)

			if err != nil {
				logger.Error(err)
			}

			cates := models.GetAllCates()
			links := models.GetAllLinks()
			archives, err := models.PostArchive()
			if err != nil {
				logger.Error(err)
			}

			c.SharedData["Moods"] = mood
			c.SharedData["Cates"] = cates
			c.SharedData["Links"] = links
			c.SharedData["Archives"] = archives

			comments, err := models.NewComments()
			if err != nil {
				logger.Error(err)
			}

			c.SharedData["NewComments"] = comments

			c.SharedData["IsAdminPage"] = false
		} else {
			c.SharedData["IsAdminPage"] = true
		}

	}

	c.Next()

	return nil
}
