package models

import "time"

type Links struct {
	Id        uint 		`form:"id" xorm:"pk"`
	Name      string    `form:"name" xorm:"varchar(100) unique"`
	Url       string    `form:"url" xorm:"varchar(200)"`
	Desc      string    `form:"desc" xorm:"varchar(255) notnull"`
	CreatedAt time.Time `form:"-" xorm:"created notnull"`
}

func (l *Links) Get() (*Links, error) {
	_, err := orm.Get(l)
	return l, err
}

func (l *Links) GetList(start int, num int) ([]*Links, error) {
	var links = make([]*Links, 0)
	start = (start - 1) * num
	err := orm.Limit(num, start).Find(&links)
	return links, err
}

func (l *Links) All() ([]*Links, error) {
	var posts = make([]*Links, 0)
	err := orm.Find(&posts)
	return posts, err
}

func (l *Links) Insert() (int64, error) {
	affected, err := orm.Insert(l)
	return affected, err
}

func (l *Links) Update() (int64, error) {
	affected, err := orm.Id(l.Id).Update(l)
	return affected, err
}

func (l *Links) Count() (int64, error) {
	affected, err := orm.Count(l)
	return affected, err
}

func (l *Links) Delete() (int64, error) {
	affected, err := orm.Id(l.Id).Delete(l)
	return affected, err
}
