package core

import (
	"time"

	"github.com/fifsky/goblog/core/ding"
	"github.com/fifsky/goblog/models"
	"github.com/ilibs/gosql"
	"github.com/ilibs/logger"
)

func StartCron() {
	t := time.NewTicker(60 * time.Second)

	for {
		select {
		case t1 := <-t.C:
			dingRemind(t1)
		}
	}
}

func dingRemind(t time.Time) {
	reminds := make([]*models.Reminds, 0)
	err := gosql.Model(&reminds).All()
	if err != nil {
		logger.Error(err)
	}

	//天气预报
	//go helpers.SendWeather(t)
	//每日一文
	//go helpers.SaveMeiRiYiWen(t)
	//每日一句
	//go helpers.SaveMeiRiYiJu(t)

	for _, v := range reminds {
		remind_date := v.RemindDate
		//fmt.Println(v, t.Format("2006-01-02 15:04:00"), remind_date.Format("2006-01-02 15:04:00"))

		content := "提醒时间:" + remind_date.Format("2006-01-02 15:04:00") + "\n提醒内容:" + v.Content

		switch v.Type {
		case 0: //固定时间
			if t.Format("2006-01-02 15:04:00") == remind_date.Format("2006-01-02 15:04:00") {
				ding.Alarm(content)
			}
		case 1: //每分钟
			ding.Alarm(content)
		case 2: //每小时
			if t.Format("04:00") == remind_date.Format("04:00") {
				ding.Alarm(content)
			}
		case 3: //每天
			if t.Format("15:04:00") == remind_date.Format("15:04:00") {
				ding.Alarm(content)
			}
		case 4: //每周
			if t.Weekday().String() == remind_date.Weekday().String() && t.Format("15:04:00") == remind_date.Format("15:04:00") {
				ding.Alarm(content)
			}
		case 5: //每月
			if t.Format("02 15:04:00") == remind_date.Format("02 15:04:00") {
				ding.Alarm(content)
			}
		case 6: //每年
			if t.Format("01-02 15:04:00") == remind_date.Format("01-02 15:04:00") {
				ding.Alarm(content)
			}
		}
	}
}