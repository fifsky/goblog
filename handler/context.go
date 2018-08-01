package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/helpers"
)

type Context struct {
	*gin.Context
}

type Response struct {
	Retcode int         `json:"retcode"`
	Data    interface{} `json:"data"`
	Msg     string      `json:"msg"`
}

func (c *Context) Fail(code int, msg interface{}) *Response {
	var message string
	if m, ok := msg.(error); ok {
		message = m.Error()
	} else {
		message = helpers.ToStr(msg)
	}

	return &Response{
		Retcode: code,
		Msg:     message,
	}
}

func (c *Context) Success(data interface{}) *Response {
	return &Response{
		Retcode: 0,
		Data:    data,
	}
}

type HandlerFunc func(c *Context) *Response

func Wrap(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := handler(&Context{c})
		c.JSON(200, response)
	}
}
