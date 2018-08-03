package core

import (
	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/helpers"
)

type Context struct {
	*gin.Context
}

func (c *Context) Fail(code int, msg interface{}) Response {
	var message string
	if m, ok := msg.(error); ok {
		message = m.Error()
	} else {
		message = helpers.ToStr(msg)
	}

	return &JSONResponse{
		Context:    c.Context,
		StatusCode: code,
		Message:    message,
	}
}

func (c *Context) Success(data interface{}) Response {
	return &JSONResponse{
		Context:    c.Context,
		StatusCode: 0,
		Data:       data,
	}
}

func (c *Context) Redirect(code int, location string) Response {
	return &RedirectResponse{
		Context:  c.Context,
		Code:     code,
		Location: location,
	}
}

func (c *Context) String(code int, name string, obj interface{}) Response {
	return &StringResponse{
		Context: c.Context,
		Name:    name,
		Data:    obj,
	}
}

func (c *Context) HTML(code int, name string, obj interface{}) Response {
	return &HTMLResponse{
		Context: c.Context,
		Name:    name,
		Data:    obj,
	}
}
