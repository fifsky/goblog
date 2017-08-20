package models

import "time"

// table users
type Users struct {
	Id        uint `xorm:"pk"`
	Name      string    `xorm:"varchar(100) unique"`  //用户名
	Password  string    `xorm:"varchar(100) notnull"` //密码
	Email     string    `xorm:"varchar(100) notnull"` //邮箱
	Type      uint8     `xorm:"notnull"`              //1:管理员,2:编辑
	Status    uint8     `xorm:"notnull"`              //1正常，2删除
	NickName  string    `xorm:"varchar(100) notnull"` // 昵称
	CreatedAt time.Time `xorm:"created notnull"`
	UpdatedAt time.Time `xorm:"updated notnull"`
}

func (u *Users) Get() (*Users, error) {
	_, err := orm.Get(u)
	return u, err
}

func (u *Users) All() ([]*Cates, error) {
	var posts = make([]*Cates, 0)
	err := orm.Find(&posts)
	return posts, err
}

func (u *Users) Insert() (int64, error) {
	affected, err := orm.Insert(u)
	return affected, err
}

func (u *Users) Update() (int64, error) {
	affected, err := orm.Id(u.Id).Update(u)
	return affected, err
}
