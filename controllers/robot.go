package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"io/ioutil"
	"github.com/fifsky/goblog/helpers"
	"github.com/tidwall/gjson"
)

type BearyChatUrl struct {
	Url string `json:"text"`
}

type BearyChatMessage struct {
	Title  string         `json:"title"`
	Text   string         `json:"text"`
	Color  string         `json:"color"`
	Images []BearyChatUrl `json:"images,omitempty"`
}

type BearyChatResponse struct {
	Text        string             `json:"text"`
	Attachments []BearyChatMessage `json:"attachments,omitempty"`
}

func BearyChat(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.String(http.StatusOK, helpers.JsonEncode(&BearyChatResponse{Text: err.Error()}))
	}

	req := gjson.ParseBytes(body)

	fmt.Println(req.Get("token"))

	msg := make([]BearyChatMessage, 0)
	msg = append(msg, BearyChatMessage{
		Title: "test",
		Text:  "test11111",
		Color: "#5788D9",
	})

	resp := &BearyChatResponse{
		Text:        req.Get("text").String(),
		Attachments: msg,
	}

	c.String(http.StatusOK, helpers.JsonEncode(resp))
}
