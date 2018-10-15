package handler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/fifsky/goblog/config"
	"github.com/fifsky/goblog/core"
	"github.com/tidwall/gjson"
)

var Handle404 core.HandlerFunc = func(c *core.Context) core.Response {
	return c.Message("未找到(404 Not Found)", "抱歉，您浏览的页面未找到。")
}

func TCaptchaVerify(ticket, randstr, ip string) error {
	p := &url.Values{}
	p.Add("aid", config.App.Common.TCaptchaId)
	p.Add("AppSecretKey", config.App.Common.TCaptchaSecret)
	p.Add("Ticket", ticket)
	p.Add("Randstr", randstr)
	p.Add("UserIP", ip)

	fmt.Println("https://ssl.captcha.qq.com/ticket/verify?"+ p.Encode())

	req, err := http.Get("https://ssl.captcha.qq.com/ticket/verify?"+ p.Encode())

	if err != nil {
		return err
	}
	defer req.Body.Close()

	str, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	ret := gjson.ParseBytes(str)
	if ret.Get("response").Int() != 1 {
		return errors.New(ret.Get("err_msg").String())
	}

	return nil
}
