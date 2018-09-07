package core

import (
	"github.com/gin-gonic/gin"
)

type IHandler interface {
	Handle(c *Context) Response
}

type HandlerFunc func(c *Context) Response

func (h HandlerFunc) Handle(c *Context) Response {
	return h(c)
}

const contextKey = "__context"

func Handle(handler IHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := handler.Handle(getContext(c))
		if response != nil {
			response.Render()
		}
	}
}

func Middlerware(handler IHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := handler.Handle(getContext(c))
		if response != nil {
			c.Abort()
			response.Render()
		}
	}
}

func getContext(c *gin.Context) *Context {
	ctx, ok := c.Get(contextKey)
	var ctx1 *Context
	if !ok {
		ctx1 = &Context{
			Context: c,
		}
		c.Set(contextKey, ctx1)
	} else {
		ctx1 = ctx.(*Context)
	}
	return ctx1
}
