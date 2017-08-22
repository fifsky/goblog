package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/models"
	"github.com/sirupsen/logrus"
	"strings"
)

func defaultH(c *gin.Context) gin.H {
	user, _ := c.Get("LoginUser")
	options := c.GetStringMapString("options")

	h := gin.H{
		"SiteTitle":   options["site_title"],
		"SiteKeyword": options["site_keyword"],
		"SiteDesc":    options["site_desc"],
		"LoginUser":   user,
	}

	url := strings.Split(c.Request.URL.Path, "/")
	h["UrlPath"] = c.Request.URL.Path

	if url[1] != "admin" {
		moodModel := new(models.Moods)
		mood, err := moodModel.Frist()
		if err != nil {
			logrus.Error(err)
		}

		cateModel := new(models.Cates)
		cates, err := cateModel.All()
		if err != nil {
			logrus.Error(err)
		}

		linkModel := new(models.Links)
		links, err := linkModel.All()
		if err != nil {
			logrus.Error(err)
		}

		h["Mood"] = mood
		h["Cates"] = cates
		h["Links"] = links

		h["IsAdminPage"] = false
	} else {
		h["IsAdminPage"] = true
	}

	return h
}
