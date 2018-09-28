package core

import (
	"bytes"
	"html/template"

	"github.com/fifsky/goblog/helpers"
	"github.com/fifsky/goblog/helpers/pagination"
	"github.com/gin-gonic/gin"
	"github.com/ilibs/sessions"
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
	SharedData gin.H
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
		StatusCode: 200,
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

func (c *Context) HTML(name string, objs ...gin.H) Response {
	data := make(map[string]interface{})
	for k, v := range c.SharedData {
		data[k] = v
	}

	for _, obj := range objs {
		for k, v := range obj {
			data[k] = v
		}
	}

	return &HTMLResponse{
		HttpStatus: getHttpStatus(c, 200),
		Context:    c.Context,
		Name:       name,
		Data:       data,
	}
}

func (c *Context) Message(title, msg string) Response {
	return c.HTML("error/message", gin.H{
		"Title":   title,
		"Message": msg,
	})
}

func (c *Context) ErrorMessage(err error) Response {
	return c.HTML("error/message", gin.H{
		"Title":   "系统错误",
		"Message": err.Error(),
	})
}

var templ = template.Must(template.New("").Funcs(funcMap).ParseGlob("views/**/*"))

func (c *Context) HTMLRender(name string, obj interface{}) (string, error) {
	writer := &bytes.Buffer{}
	err := templ.ExecuteTemplate(writer, name, obj)
	if err != nil {
		return "", err
	}
	return writer.String(), nil
}
