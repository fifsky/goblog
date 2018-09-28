package models

import (
	"strconv"
	"time"

	"github.com/ilibs/gosql"
)

type Users struct {
	Id        int       `form:"id" json:"id" db:"id"`
	Name      string    `form:"name" json:"name" db:"name"`
	Password  string    `form:"password" json:"password" db:"password"`
	NickName  string    `form:"nick_name" json:"nick_name" db:"nick_name"`
	Email     string    `form:"email" json:"email" db:"email"`
	Status    int       `form:"status" json:"status" db:"status"`
	Type      int       `form:"type" json:"type" db:"type"`
	CreatedAt time.Time `form:"-" json:"created_at" db:"created_at"`
	UpdatedAt time.Time `form:"-" json:"updated_at" db:"updated_at"`
}

func (u *Users) DbName() string {
	return "default"
}

func (u *Users) TableName() string {
	return "users"
}

func (u *Users) PK() string {
	return "id"
}

func (u *Users) AfterChange() {
	Cache.Delete("user:" + strconv.Itoa(u.Id))
}

func GetUser(uid int) (*Users, error) {
	if v, ok := Cache.Get("user:" + strconv.Itoa(uid)); ok {
		return v.(*Users), nil
	}

	user := &Users{}
	err := gosql.Model(user).Where("id = ?", uid).Get()

	if err != nil {
		return nil, err
	}
	Cache.Set("user:"+strconv.Itoa(uid), user, 1*time.Hour)

	return user, nil
}

func UserGetList(start int, num int) ([]*Users, error) {
	var m = make([]*Users, 0)
	start = (start - 1) * num
	err := gosql.Model(&m).OrderBy("id desc").Limit(num).Offset(start).All()
	if err != nil {
		return nil, err
	}
	return m, nil
}
