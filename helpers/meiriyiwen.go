package helpers

import (
	"net/http"
	"io/ioutil"
	"github.com/tidwall/gjson"
	"fmt"
	"time"
	"github.com/fifsky/goblog/models"
)

func SaveMeiRiYiWen(t time.Time) {
	if t.Format("15:04") != "17:40" {
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
		m.Content = fmt.Sprintf("<p>文/%s</p>%s", ret.Get("data.author"), ret.Get("data.content").String())
		m.UserId = 1
		m.Title = ret.Get("data.title").String()

		m.Insert()
	}
}
