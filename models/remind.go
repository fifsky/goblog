package models

import (
	"time"
	"github.com/sirupsen/logrus"
)

type Reminds struct {
	Id         uint      `form:"id" xorm:"pk"`
	Type       int       `form:"type" xorm:"int"`
	At         string    `form:"at" xorm:"varchar(20)"`
	Content    string    `form:"content" xorm:"varchar(255)"`
	RemindDate string    `form:"remind_date" xorm:"varchar(50)"`
	CreatedAt  time.Time `form:"-" xorm:"created notnull"`
}

func (l *Reminds) Get() (*Reminds, bool) {
	has, err := orm.Get(l)
	if err != nil {
		logrus.Error(err)
		return l, false
	}
	return l, has
}

func (l *Reminds) GetList(start int, num int) ([]*Reminds, error) {
	var links = make([]*Reminds, 0)
	start = (start - 1) * num
	err := orm.Limit(num, start).Find(&links)
	return links, err
}

func (l *Reminds) All() ([]*Reminds, error) {
	var posts = make([]*Reminds, 0)
	err := orm.Find(&posts)
	return posts, err
}

func (l *Reminds) Insert() (int64, error) {
	affected, err := orm.Insert(l)
	return affected, err
}

func (l *Reminds) Update() (int64, error) {
	affected, err := orm.Id(l.Id).Update(l)
	return affected, err
}

func (l *Reminds) Count() (int64, error) {
	affected, err := orm.Count(l)
	return affected, err
}

func (l *Reminds) Delete() (int64, error) {
	affected, err := orm.Id(l.Id).Delete(l)
	return affected, err
}
