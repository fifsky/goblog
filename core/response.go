package core

import (
	"github.com/gin-gonic/gin"
)

type Response interface {
	Render()
}

type JSONResponse struct {
	Context    *gin.Context `json:"-"`
	StatusCode int          `json:"statusCode"`
	Data       interface{}  `json:"data"`
	Message    string       `json:"message"`
}

func (c *JSONResponse) Render() {
	c.Context.JSON(200, c)
}

type RedirectResponse struct {
	Context  *gin.Context
	Code     int
	Location string
}

func (c *RedirectResponse) Render() {
	c.Context.Redirect(c.Code, c.Location)
}

type StringResponse struct {
	Context *gin.Context
	Name    string
	Data    interface{}
}

func (c *StringResponse) Render() {
	c.Context.String(200, c.Name, c.Data)
}

type HTMLResponse struct {
	Context *gin.Context
	Name    string
	Data    interface{}
}

func (c *HTMLResponse) Render() {
	c.Context.HTML(200, c.Name, c.Data)
}
