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
	CreatedAt time.Time `xorm:"notnull"`
	UpdatedAt time.Time `xorm:"notnull"`
	DeletedAt *time.Time `xorm:"notnull"`
}
