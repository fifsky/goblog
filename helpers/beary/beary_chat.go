package beary

import (
	"encoding/json"
	"io/ioutil"
	"github.com/sirupsen/logrus"
	"net/http"
	"bytes"
	"time"
	"github.com/pkg/errors"
	"fmt"
)

var TOKEN string

type Url struct {
	Url string `json:"text"`
}

type Message struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	Url    string `json:"url,omitempty"`
	Color  string `json:"color"`
	Images []Url  `json:"images,omitempty"`
}

type Request struct {
	Text        string    `json:"text"`
	Markdown    bool      `json:"markdown,omitempty"`
	Channel     string    `json:"channel,omitempty"`
	User        string    `json:"user,omitempty"`
	Attachments []Message `json:"attachments,omitempty"`
}

type Response struct {
	Code   int
	Error  string
	Result string
}

func Alarm(content string, channel string, at string) error {
	req := &Request{
		Text:    content,
		Channel: channel,
		User:    at,
	}

	return Send(req)
}

func Send(req *Request) error {
	buf, err := json.Marshal(req)
	if err != nil {
		return err
	}

	url := "https://hook.bearychat.com/=bwCyz/incoming/" + TOKEN
	resp, err := PostJson(url, buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if data, err := ioutil.ReadAll(resp.Body); err == nil {
		ret := new(Response)
		err := json.Unmarshal(data, ret)
		if err != nil || ret.Code != 0 {
			logrus.Debug("bearychat", string(data))
			return errors.New(fmt.Sprintf("beary response error[%d]:%s", ret.Code, ret.Error))
		}
	}

	if err != nil {
		logrus.Debug("alarm", err.Error())
		return err
	}
	return nil
}

func PostJson(url string, data []byte) (*http.Response, error) {

	fmt.Println(string(data))

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
