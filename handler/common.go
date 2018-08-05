package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/models"
	"github.com/ilibs/gosql"
	"github.com/ilibs/logger"
	"github.com/fifsky/goblog/core"
)

func defaultH(c *gin.Context) gin.H {
	user, _ := c.Get("LoginUser")
	options := c.GetStringMapString("options")

	h := gin.H{
		"SiteTitle":   options["site_name"],
		"SiteKeyword": options["site_keyword"],
		"SiteDesc":    options["site_desc"],
		"LoginUser":   user,
	}

	url := strings.Split(c.Request.URL.Path, "/")
	h["UrlPath"] = c.Request.URL.Path
	h["URI"] = c.Request.RequestURI

	if url[1] != "admin" {
		mood, err := models.MoodFrist()

		if err != nil {
			logger.Error(err)
		}

		cates := make([]*models.Cates, 0)
		err = gosql.Model(&cates).All()
		if err != nil {
			logger.Error(err)
		}

		links := make([]*models.Links, 0)
		err = gosql.Model(&links).All()
		if err != nil {
			logger.Error(err)
		}

		archives, err := models.PostArchive()
		if err != nil {
			logger.Error(err)
		}

		h["Mood"] = mood
		h["Cates"] = cates
		h["Links"] = links
		h["Archives"] = archives

		h["IsAdminPage"] = false
	} else {
		h["IsAdminPage"] = true
	}

	return h
}

var Handle404 core.HandlerFunc = func(c *core.Context) core.Response {
	return HandleMessage(c, "未找到(404 Not Found)", "抱歉，您浏览的页面未找到。")
}

func HandleMessage(c *core.Context, title string, message string) core.Response {
	return c.HTML("error/message", gin.H{
		"Title":   title,
		"Message": message,
	})
}
