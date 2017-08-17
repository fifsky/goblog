package models

type Cates struct {
	BaseModel
	Name   string    `xorm:"varchar(100) unique 'un_name'"` //用户名
	Desc   string    `xorm:"varchar(255) notnull"`     //密码
	Domain string    `xorm:"varchar(100) "`             //邮箱
}