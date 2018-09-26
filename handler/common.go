package handler

import (
	"runtime"
	"strings"

	"github.com/fifsky/goblog/core"
	"github.com/fifsky/goblog/models"
	"github.com/gin-gonic/gin"
	"github.com/ilibs/logger"
)

func defaultH(c *gin.Context) gin.H {
	user, _ := c.Get("LoginUser")
	options := c.GetStringMapString("options")

	h := gin.H{
		"SiteTitle":   options["site_name"],
		"SiteKeyword": options["site_keyword"],
		"SiteDesc":    options["site_desc"],
		"LoginUser":   user,
		"GOVERSION":   runtime.Version(),
	}

	url := strings.Split(c.Request.URL.Path, "/")
	h["UrlPath"] = c.Request.URL.Path
	h["URI"] = c.Request.RequestURI

	if url[1] != "admin" {
		mood, err := models.MoodFrist()

		if err != nil {
			logger.Error(err)
		}

		cates := models.GetAllCates()
		links := models.GetAllLinks()
		archives, err := models.PostArchive()
		if err != nil {
			logger.Error(err)
		}

		h["Mood"] = mood
		h["Cates"] = cates
		h["Links"] = links
		h["Archives"] = archives

		comments,err := models.NewComments()
		if err != nil {
			logger.Error(err)
		}

		h["NewComments"] = comments

		h["IsAdminPage"] = false
	} else {
		h["IsAdminPage"] = true
	}

	return h
}

var Handle404 core.HandlerFunc = func(c *core.Context) core.Response {
	return c.Message("未找到(404 Not Found)", "抱歉，您浏览的页面未找到。")
}