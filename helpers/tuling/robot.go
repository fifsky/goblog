package tuling

import (
	"encoding/json"
	"io/ioutil"
	"github.com/sirupsen/logrus"
	"net/http"
	"bytes"
	"time"
	"github.com/tidwall/gjson"
)

var TOKEN string

type InputText struct {
	Text string `json:"text"`
}

type InputImage struct {
	Url string `json:"url"`
}

type UserInfo struct {
	ApiKey string `json:"apiKey"`
	UserId string `json:"userId"`
}

type Location struct {
	City     string `json:"city"`
	Province string `json:"province"`
	Street   string `json:"street"`
}

type SelfInfo struct {
	Location *Location `json:"location"`
}

type Perception struct {
	InputText  *InputText  `json:"inputText"`
	InputImage *InputImage `json:"inputImage,omitempty"`
	SelfInfo   *SelfInfo   `json:"selfInfo,omitempty"`
}

type Request struct {
	ReqType    int         `json:"reqType"`
	Perception *Perception `json:"perception,omitempty"`
	UserInfo   *UserInfo   `json:"userInfo,omitempty"`
}

func Say(content string, image string) (map[string]string, error) {
	per := &Perception{
		InputText: &InputText{
			Text: content,
		},
		SelfInfo: &SelfInfo{
			Location: &Location{
				City:     "上海",
				Province: "上海",
				Street:   "九新公路",
			},
		},
	}

	user := &UserInfo{
		ApiKey: TOKEN,
		UserId: "10001",
	}

	req := &Request{
		ReqType:    0,
		Perception: per,
		UserInfo:   user,
	}

	if image != "" {
		req.ReqType = 1
		req.Perception.InputImage = &InputImage{
			Url: image,
		}
	}

	return send(req)
}

func send(req *Request) (map[string]string, error) {
	buf, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := "http://openapi.tuling123.com/openapi/api/v2"
	resp, err := postJson(url, buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Debug("alarm", err.Error())
		return nil, err
	}

	ret := gjson.ParseBytes(data)
	//fmt.Println(ret.Get("results").Array())

	m := make(map[string]string, 0)
	//文本(text);连接(url);音频(voice);视频(video);图片(image);图文(news)
	for _, v := range ret.Get("results").Array() {
		t := v.Get("resultType").String()
		m[t] = v.Get("values." + t).String()
	}

	return m, nil
}

func postJson(url string, data []byte) (*http.Response, error) {
	//fmt.Println(string(data))

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
