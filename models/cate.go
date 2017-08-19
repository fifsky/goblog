package models

import "time"

type Cates struct {
	Id        uint `xorm:"pk"`
	Name      string    `xorm:"varchar(100) unique"` //用户名
	Desc      string    `xorm:"varchar(255) notnull"`          //密码
	Domain    string    `xorm:"varchar(100) "`                 //邮箱
	CreatedAt time.Time `xorm:"notnull"`
	UpdatedAt time.Time `xorm:"notnull"`
	DeletedAt *time.Time `xorm:"notnull"`
}


func (c *Cates) Get() (*Cates, error) {
	_, err := orm.Get(c)
	return c, err
}