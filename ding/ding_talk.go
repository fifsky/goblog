package ding

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var DING_TALK_TOKEN string

type DingTalkRequest struct {
	Msgtype string                 `json:"msgtype"`
	Text    map[string]string      `json:"text"`
	At      map[string]interface{} `json:"at"`
}

type DingTalkResponse struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func Alarm(content string, at ...string) error {
	dingtalk := &DingTalkRequest{
		Msgtype: "text",
		Text: map[string]string{
			"content": content,
		},
	}

	if len(at) > 0 {
		dingtalk.At = map[string]interface{}{
			"atMobiles": at,
		}
	}

	buf, err := json.Marshal(dingtalk)
	if err != nil {
		return err
	}

	url := "https://oapi.dingtalk.com/robot/send?access_token=" + DING_TALK_TOKEN
	resp, err := PostJson(url, buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if data, err := ioutil.ReadAll(resp.Body); err == nil {
		ret := &DingTalkResponse{}
		err := json.Unmarshal(data, ret)
		if err != nil {
			return err
		}

		if ret.Errcode != 0 {
			return errors.New("ding response error:" + ret.Errmsg + "[" + strconv.Itoa(ret.Errcode) + "]")
		}
	}

	if err != nil {
		return err
	}
	return nil
}

func PostJson(url string, data []byte) (*http.Response, error) {
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	client := &http.Client{}
	client.Timeout = 5 * time.Second

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, err
}
