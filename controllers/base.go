package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/models"
	"github.com/sirupsen/logrus"
	"strings"
)

func defaultH(c *gin.Context) gin.H {
	user, _ := c.Get("LoginUser")
	options := c.MustGet("options").(map[string]string)

	h := gin.H{
		"SiteTitle": options["site_name"],
		"LoginUser": user,
	}

	url := strings.Split(c.Request.URL.Path, "/")
	h["UrlPath"] = c.Request.URL.Path

	if url[1] != "admin" {
		moodModel := new(models.Moods)
		mood, err := moodModel.Frist()
		if err != nil {
			logrus.Error("get mood error:" + err.Error())
		}
		h["Mood"] = mood
		h["IsAdminPage"] = false
	} else {
		h["IsAdminPage"] = true
	}

	return h
}
