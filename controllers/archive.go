package controllers

import (
	"github.com/gin-gonic/gin"
	//"github.com/fifsky/goblog/models"
	//"net/http"
	//"strconv"
)

func ArchiveGet(c *gin.Context) {
	//year := c.Param("year")
	//month := c.Param("month")
	//posts, err := models.ListPostByArchive(year, month)
	//if err == nil {
	//	for _, post := range posts {
	//		post.Tags, _ = models.ListTagByPostId(strconv.FormatUint(uint64(post.ID), 10))
	//	}
	//	c.HTML(http.StatusOK, "index/index.html", gin.H{
	//		"posts":    posts,
	//		"tags":     models.MustListTag(),
	//		"archives": models.MustListPostArchives(),
	//		"links":    models.MustListLinks(),
	//	})
	//} else {
	//	c.AbortWithStatus(http.StatusInternalServerError)
	//}
}
