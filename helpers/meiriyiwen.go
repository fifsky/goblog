package helpers

import (
	"net/http"
	"io/ioutil"
	"github.com/tidwall/gjson"
	"fmt"
	"time"
	"github.com/fifsky/goblog/models"
	"github.com/fifsky/goblog/helpers/beary"
)

func SaveMeiRiYiWen(t time.Time) {
	if t.Format("15:04") != "10:00" {
		return
	}

	resp, err := http.Get("https://interface.meiriyiwen.com/article/today?dev=1")
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if data, err := ioutil.ReadAll(resp.Body); err == nil {
		ret := gjson.ParseBytes(data)

		m := &models.Posts{}
		m.CateId = 5
		m.Type = 1
		m.Content = fmt.Sprintf(`<p style="text-align: center;">文/%s</p>%s`, ret.Get("data.author"), ret.Get("data.content").String())
		m.UserId = 1
		m.Title = ret.Get("data.title").String()

		m.Insert()

		msg := make([]beary.Message, 0)
		msg = append(msg, beary.Message{
			Title: m.Title,
			Url:   fmt.Sprintf("https://fifsky.com/article/%d", m.Id),
			Color: "#1FBECA",
		})

		req := &beary.Request{
			Text:        "发表了每日一文",
			Channel:     "豆爸的私人助理",
			Attachments: msg,
		}
		beary.Send(req)
	}

}

func SaveMeiRiYiJu(t time.Time) {
	if t.Format("15:04") != "09:00" {
		return
	}

	resp, err := http.Get("http://open.iciba.com/dsapi/")
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if data, err := ioutil.ReadAll(resp.Body); err == nil {
		ret := gjson.ParseBytes(data)

		m := &models.Moods{}
		m.Content = ret.Get("content").String() + "<br>" + ret.Get("note").String()
		m.UserId = 1
		m.Insert()

		msg := make([]beary.Message, 0)
		msg = append(msg, beary.Message{
			Text:  m.Content,
			Color: "#95B9B3",
		})

		req := &beary.Request{
			Text:        "发表了每日一句",
			Channel:     "豆爸的私人助理",
			Attachments: msg,
		}
		beary.Send(req)
	}
}
