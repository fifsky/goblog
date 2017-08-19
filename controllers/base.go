package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/models"
	"github.com/sirupsen/logrus"
)

func defaultH(c *gin.Context) gin.H {
	options := c.MustGet("options").(map[string]string)

	moodModel := new(models.Moods)
	mood, err := moodModel.Frist()

	if err != nil {
		logrus.Error("get mood error:" + err.Error())
	}

	return gin.H{
		"SiteTitle": options["site_name"],
		"Mood":      mood,
	}
}
