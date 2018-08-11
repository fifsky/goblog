package handler

import (
	"github.com/fifsky/goblog/core"
	"github.com/fifsky/goblog/models"
	"github.com/ilibs/gosql"
	"github.com/ilibs/logger"
	"time"
	"github.com/gin-gonic/gin"
	"fmt"
)

var CommentPost core.HandlerFunc = func(c *core.Context) core.Response {
	comment := &models.Comments{}
	if err := c.ShouldBind(comment); err != nil {
		return c.Fail(201, "参数错误:"+err.Error())
	}

	if comment.Name == "" {
		return c.Fail(201, "昵称不能为空")
	}

	if comment.Content == "" {
		return c.Fail(201, "评论内容不能为空")
	}

	if comment.PostId <= 0 {
		return c.Fail(201, "非法评论")
	}

	comment.CreatedAt = time.Now()
	comment.IP = c.ClientIP()

	if _, err := gosql.Model(comment).Create(); err != nil {
		logger.Error(err)
		return c.Fail(201, "评论失败"+err.Error())
	}

	body, err := c.HTMLRender("layout/comment_item", comment)
	if err != nil {
		return c.Fail(202, err)
	}

	fmt.Println(body)

	return c.Success(gin.H{
		"content": body,
	})
}
