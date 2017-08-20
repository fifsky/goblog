package models

import "time"

type Cates struct {
	Id        uint `xorm:"pk"`
	Name      string    `xorm:"varchar(100) unique"`  //分类名
	Desc      string    `xorm:"varchar(255) notnull"` //分类详情
	Domain    string    `xorm:"varchar(100) unique"`  //分类短域名
	CreatedAt time.Time `xorm:"created notnull"`
	UpdatedAt time.Time `xorm:"updated notnull"`
}

func (c *Cates) Get() (*Cates, error) {
	_, err := orm.Get(c)
	return c, err
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
