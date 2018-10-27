package core

import (
	"github.com/gin-gonic/gin"
)

type Response interface {
	Render()
}

type JSONResponse struct {
	HttpStatus int          `json:"-"`
	Context    *gin.Context `json:"-"`
	Data       interface{}  `json:"data"`
}

func (c *JSONResponse) Render() {
	c.Context.JSON(c.HttpStatus, c.Data)
}

type XMLResponse struct {
	HttpStatus int          `json:"-"`
	Context    *gin.Context `json:"-"`
	Data       interface{}  `json:"data"`
}

func (c *XMLResponse) Render() {
	c.Context.XML(c.HttpStatus, c.Data)
}

type ApiResponse struct {
	HttpStatus int          `json:"-"`
	Context    *gin.Context `json:"-"`
	StatusCode int          `json:"statusCode"`
	Data       interface{}  `json:"data"`
	Message    string       `json:"message"`
}

func (c *ApiResponse) Render() {
	c.Context.JSON(c.HttpStatus, c)
}

type RedirectResponse struct {
	HttpStatus int          `json:"-"`
	Context    *gin.Context `json:"-"`
	Location   string
}

func (c *RedirectResponse) Render() {
	c.Context.Redirect(c.HttpStatus, c.Location)
}

type StringResponse struct {
	HttpStatus int          `json:"-"`
	Context    *gin.Context `json:"-"`
	Name       string
	Data       []interface{}
}

func (c *StringResponse) Render() {
	c.Context.String(c.HttpStatus, c.Name, c.Data...)
}

type HTMLResponse struct {
	HttpStatus int          `json:"-"`
	Context    *gin.Context `json:"-"`
	Name       string
	Data       interface{}
}

func (c *HTMLResponse) Render() {
	c.Context.HTML(c.HttpStatus, c.Name, c.Data)
}