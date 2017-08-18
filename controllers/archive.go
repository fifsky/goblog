package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/fifsky/goblog/models"
	"github.com/fifsky/goblog/helpers"
)

func ArchiveGet(c *gin.Context) {
	id, _ := helpers.StrTo(c.Param("id")).Uint()
	postModel := &models.Posts{Id: id}
	post, err := postModel.Get()

	cateModel := &models.Cates{Id: post.CateId}
	cate, err := cateModel.Get()

	userModel := &models.Users{Id: post.UserId}
	user, err := userModel.Get()

	newpost := &models.UserPosts{Posts: *post, Name: cate.Name, NickName: user.NickName}

	moodModel := new(models.Moods)
	mood, err := moodModel.Frist()

	h := defaultH(c)
	h["Title"] = post.Title
	h["Post"] = newpost
	h["Mood"] = mood

	if err == nil {
		c.HTML(http.StatusOK, "index/article", h)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
