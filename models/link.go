package models

import "time"

type Links struct {
	Id        uint `xorm:"pk"`
	Name      string    `xorm:"varchar(100) unique"`
	Url       string    `xorm:"varchar(200)"`
	Desc      string    `xorm:"varchar(255) notnull"`
	CreatedAt time.Time `xorm:"created notnull"`
	UpdatedAt time.Time `xorm:"updated notnull"`
}

func (l *Links) Get() (*Links, error) {
	_, err := orm.Get(l)
	return l, err
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
