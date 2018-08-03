package core

import (
	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/context"
)

type IHandler interface {
	Handle(c *context.Context) Response
}

type HandlerFunc func(c *context.Context) Response

func (h HandlerFunc) Handle(c *context.Context) Response {
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

func Middware(handler IHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := handler.Handle(getContext(c))
		if response != nil {
			c.Abort()
			response.Render()
		}
	}
}

func getContext(c *gin.Context) *context.Context {
	ctx, ok := c.Get(contextKey)
	var ctx1 *context.Context
	if !ok {
		ctx1 = &context.Context{
			Context: c,
		}
		c.Set(contextKey, ctx1)
	} else {
		ctx1 = ctx.(*context.Context)
	}
	return ctx1
}
