package handler

import (
	"fmt"
	"time"

	"github.com/fifsky/goblog/core"
	"github.com/fifsky/goblog/ding"
	"github.com/fifsky/goblog/models"
	"github.com/gin-gonic/gin"
	"github.com/ilibs/gosql"
	"github.com/ilibs/logger"
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

	if err := TCaptchaVerify(c.PostForm("ticket"), c.PostForm("randstr"), c.ClientIP()); err != nil {
		return c.Fail(201, err)
	}

	//if !captcha.VerifyString(c.PostForm("captcha_id"), c.PostForm("captcha")) {
	//	return c.Fail(201, "验证码错误")
	//}

	post := &models.Posts{}
	err := gosql.Model(post).Where("id = ?", comment.PostId).Get()
	if err != nil {
		return c.Fail(201, "文章不存在")
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

	content := "您有新的评论!\n"
	content += fmt.Sprintf("文章:%s\n", post.Title)
	content += fmt.Sprintf("评论内容:%s\n", comment.Content)
	content += fmt.Sprintf("评论昵称:%s\n", comment.Name)
	content += fmt.Sprintf("评论时间:%s\n", comment.CreatedAt.Format("2006-01-02 15:04:05"))
	content += fmt.Sprintf("评论IP:%s\n", comment.IP)

	ding.Alarm(content)

	return c.Success(gin.H{
		"content": body,
	})
}
