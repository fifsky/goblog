package models

// table users
type Users struct {
	BaseModel
	Name     string    `xorm:"unique 'un_user_name'"` //用户名
	Password string    `xorm:"notnull"`               //密码
	Email    string    `xorm:"notnull"`               //邮箱
	Type     uint8     `xorm:"notnull"`               //1:管理员,2:编辑
	Status   uint8     `xorm:"notnull"`               //1正常，2删除
	NickName string    `xorm:"notnull"`               // 昵称
}
