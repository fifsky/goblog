package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"io/ioutil"
	"github.com/fifsky/goblog/helpers"
	"github.com/tidwall/gjson"
	"strings"
	"github.com/fifsky/goblog/helpers/tuling"
)

type BearyChatUrl struct {
	Url string `json:"url"`
}

type BearyChatMessage struct {
	Title  string          `json:"title"`
	Text   string          `json:"text"`
	Color  string          `json:"color"`
	Images []*BearyChatUrl `json:"images,omitempty"`
}

type BearyChatResponse struct {
	Text        string              `json:"text"`
	Attachments []*BearyChatMessage `json:"attachments,omitempty"`
}

func BearyChat(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.String(http.StatusOK, helpers.JsonEncode(&BearyChatResponse{Text: err.Error()}))
	}

	req := gjson.ParseBytes(body)

	//fmt.Println(req.Get("token"))

	text := req.Get("text").String()
	text = strings.TrimPrefix(text, req.Get("trigger_word").String()+" ")

	//文本(text);连接(url);音频(voice);视频(video);图片(image);图文(news)
	mlist, err := tuling.Say(text, "")
	content := ""

	if err != nil {
		content = err.Error()
	}

	resp := &BearyChatResponse{
		Text: content,
	}

	msg := &BearyChatMessage{}

	//fmt.Println(mlist)

	for k, v := range mlist {

		switch k {
		case "text", "url", "voice", "video", "news":
			content += v + "\n"
		case "image":
			msg.Images = append(msg.Images, &BearyChatUrl{v})
		}

		resp.Text = content
	}
	resp.Attachments = append(resp.Attachments, msg)

	c.String(http.StatusOK, helpers.JsonEncode(resp))
}
