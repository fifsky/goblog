package context

import (
	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/core"
	"github.com/fifsky/goblog/helpers"
)

type Context struct {
	*gin.Context
}

func (c *Context) Fail(code int, msg interface{}) core.Response {
	var message string
	if m, ok := msg.(error); ok {
		message = m.Error()
	} else {
		message = helpers.ToStr(msg)
	}

	return &core.JSONResponse{
		Context:    c.Context,
		StatusCode: code,
		Message:    message,
	}
}

func (c *Context) Success(data interface{}) core.Response {
	return &core.JSONResponse{
		Context:    c.Context,
		StatusCode: 0,
		Data:       data,
	}
}

func (c *Context) Redirect(code int, location string) core.Response {
	return &core.RedirectResponse{
		Context:  c.Context,
		Code:     code,
		Location: location,
	}
}

func (c *Context) String(code int, name string, obj interface{}) core.Response {
	return &core.StringResponse{
		Context: c.Context,
		Name:    name,
		Data:    obj,
	}
}

func (c *Context) HTML(code int, name string, obj interface{}) core.Response {
	return &core.HTMLResponse{
		Context: c.Context,
		Name:    name,
		Data:    obj,
	}
}
