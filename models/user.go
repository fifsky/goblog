package models

import (
	"time"
	"github.com/sirupsen/logrus"
)

// table users
type Users struct {
	Id        uint `form:"id" xorm:"pk"`
	Name      string    `form:"name" xorm:"varchar(100) unique"`       //用户名
	NickName  string    `form:"nick_name" xorm:"varchar(100) notnull"` // 昵称
	Password  string    `form:"password" xorm:"varchar(100) notnull"`  //密码
	Email     string    `form:"email" xorm:"varchar(100) notnull"`     //邮箱
	Type      uint8     `form:"type" xorm:"notnull"`                   //1:管理员,2:编辑
	Status    uint8     `form:"status" xorm:"notnull"`                 //1正常，2删除
	CreatedAt time.Time `form:"-" xorm:"created notnull"`
	UpdatedAt time.Time `form:"-" xorm:"updated notnull"`
}

func (u *Users) Get() (*Users, bool) {
	has, err := orm.Get(u)
	if err != nil {
		logrus.Error(err)
		return u, false
	}
	return u, has
}

func (u *Users) GetList(start int, num int) ([]*Users, error) {
	var users = make([]*Users, 0)
	start = (start - 1) * num
	err := orm.Limit(num, start).Find(&users)
	return users, err
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

func (u *Users) Count() (int64, error) {
	affected, err := orm.Count(u)
	return affected, err
}

func (u *Users) Delete() (int64, error) {
	affected, err := orm.Id(u.Id).Update(u)
	return affected, err
}
