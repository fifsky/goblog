package handler

import (
	"github.com/fifsky/goblog/core"
)

var Handle404 core.HandlerFunc = func(c *core.Context) core.Response {
	return c.Message("未找到(404 Not Found)", "抱歉，您浏览的页面未找到。")
}