package core

import (
	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/helpers"
	"github.com/ilibs/sessions"
	"github.com/fifsky/goblog/helpers/pagination"
)

func getHttpStatus(c *Context, status int) int {
	if c.HttpStatus == 0 {
		return status
	}
	return c.HttpStatus
}

type Context struct {
	*gin.Context
	HttpStatus int
}

func (c *Context) Pagination(total int64, num, page int) *pagination.Paginater {
	return pagination.New(int(total), num, page, 3)
}

func (c *Context) Session() sessions.Session {
	return sessions.Default(c.Context)
}

func (c *Context) Status(status int) {
	c.HttpStatus = status
}

func (c *Context) Fail(code int, msg interface{}) Response {
	var message string
	if m, ok := msg.(error); ok {
		message = m.Error()
	} else {
		message = helpers.ToStr(msg)
	}

	return &ApiResponse{
		HttpStatus: getHttpStatus(c, 200),
		Context:    c.Context,
		StatusCode: code,
		Message:    message,
	}
}

func (c *Context) Success(data interface{}) Response {
	return &ApiResponse{
		HttpStatus: getHttpStatus(c, 200),
		Context:    c.Context,
		StatusCode: 0,
		Data:       data,
		Message:    "ok",
	}
}

func (c *Context) JSON(data interface{}) Response {
	return &JSONResponse{
		HttpStatus: getHttpStatus(c, 200),
		Context:    c.Context,
		Data:       data,
	}
}

func (c *Context) Redirect(location string) Response {
	return &RedirectResponse{
		HttpStatus: getHttpStatus(c, 302),
		Context:    c.Context,
		Location:   location,
	}
}

func (c *Context) String(format string, values ...interface{}) Response {
	return &StringResponse{
		HttpStatus: getHttpStatus(c, 200),
		Context:    c.Context,
		Name:       format,
		Data:       values,
	}
}

func (c *Context) HTML(name string, obj interface{}) Response {
	return &HTMLResponse{
		HttpStatus: getHttpStatus(c, 200),
		Context:    c.Context,
		Name:       name,
		Data:       obj,
	}
}
