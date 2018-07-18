package controllers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/models"
	"github.com/sirupsen/logrus"
	"github.com/ilibs/gosql"
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

	if url[1] != "admin" {
		mood ,err := models.MoodFrist()

		if err != nil {
			logrus.Error(err)
		}

		cates := make([]*models.Cates, 0)
		err = gosql.Model(&cates).All()
		if err != nil {
			logrus.Error(err)
		}

		links :=  make([]*models.Links, 0)
		err = gosql.Model(&links).All()
		if err != nil {
			logrus.Error(err)
		}

		archives, err := models.PostArchive()
		if err != nil {
			logrus.Error(err)
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
