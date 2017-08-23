package models

import (
	"time"
	"github.com/sirupsen/logrus"
)

type Cates struct {
	Id        uint 		`form:"id" xorm:"pk"`
	Name      string    `form:"name" xorm:"varchar(100) unique"`   //分类名
	Desc      string    `form:"desc" xorm:"varchar(255) notnull"`  //分类详情
	Domain    string    `form:"domain" xorm:"varchar(100) unique"` //分类短域名
	CreatedAt time.Time `form:"-" xorm:"created notnull"`
	UpdatedAt time.Time `form:"-" xorm:"updated notnull"`
}

func (c *Cates) Get() (*Cates, bool) {
	has, err := orm.Get(c)
	if err != nil {
		logrus.Error(err)
		return c, false
	}
	return c, has
}

func (c *Cates) GetList(start int, num int) ([]*Cates, error) {
	var posts = make([]*Cates, 0)
	start = (start - 1) * num
	err := orm.Limit(num, start).Find(&posts)
	return posts, err
}

func (c *Cates) All() ([]*Cates, error) {
	var posts = make([]*Cates, 0)
	err := orm.Find(&posts)
	return posts, err
}

func (c *Cates) Insert() (int64, error) {
	affected, err := orm.Insert(c)
	return affected, err
}

func (c *Cates) Update() (int64, error) {
	affected, err := orm.Id(c.Id).Update(c)
	return affected, err
}

func (c *Cates) Count() (int64, error) {
	affected, err := orm.Count(c)
	return affected, err
}

func (c *Cates) Delete() (int64, error) {
	affected, err := orm.Id(c.Id).Delete(c)
	return affected, err
}
