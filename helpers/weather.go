package helpers

import (
	"net/http"
	"io/ioutil"
	"github.com/tidwall/gjson"
	"github.com/fifsky/goblog/core/ding"
	"fmt"
	"time"
)

func SendWeather(t time.Time) {
	if t.Format("15:04") != "08:00" {
		return
	}

	resp, err := http.Get("http://tj.nineton.cn/Heart/index/all?city=CHSH000700&language=zh-chs&unit=c&aqi=city&alarm=1&key=78928e706123c1a8f1766f062bc8676b")
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if data, err := ioutil.ReadAll(resp.Body); err == nil {

		ret := gjson.ParseBytes(data)
		city := ret.Get("weather.0.city_name")
		weather := ret.Get("weather.0.future.0")
		ding.Alarm(fmt.Sprintf("今日天气\n%s %s\n温度 %s-%s ℃\n%s", city, weather.Get("text"), weather.Get("low"), weather.Get("high"), weather.Get("wind")))
	}
}
