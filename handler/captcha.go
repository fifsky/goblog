package handler

import (
	"path"
	"strconv"

	"github.com/fifsky/goblog/core"
	"github.com/ilibs/captcha"
)

var CaptchaGet core.HandlerFunc = func(c *core.Context) core.Response {
	_, file := path.Split(c.Request.URL.Path)
	ext := path.Ext(file)
	id := file[:len(file)-len(ext)]
	if ext == "" || id == "" {
		return c.Fail(202,"id not found")
	}

	if c.Query("reload") != "" {
		captcha.Reload(id)
	}

	var imgHeight = captcha.StdHeight
	var imgWidth = captcha.StdWidth
	if c.Query("h") != "" {
		ih, err := strconv.Atoi(c.Query("h"))
		if err == nil {
			imgHeight = ih
		}
	}

	if c.Query("w") != "" {
		iw, err := strconv.Atoi(c.Query("w"))
		if err == nil {
			imgWidth = iw
		}
	}

	switch ext {
	case ".png":
		c.Header("Content-Type", "image/png")
		captcha.WriteImage(c.Writer, id, imgWidth, imgHeight)
		return nil
	default:
		return c.Fail(202,"ext not found")
	}

	return c.Fail(202,"ext not found")
}
