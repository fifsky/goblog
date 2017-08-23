package models

import (
	"time"
	"github.com/sirupsen/logrus"
)

type Moods struct {
	Id        uint        `form:"id" xorm:"pk"`
	Content   string      `form:"content" xorm:"varchar(255) notnull"`
	UserId    uint        `form:"-" xorm:"notnull"`
	CreatedAt time.Time   `form:"-" xorm:"created notnull"`
}

type UserMoods struct {
	Moods    `xorm:"extends"`
	NickName string
}

func (UserMoods) TableName() string {
	return "moods"
}

func (m *Moods) Frist() (*Moods, error) {
	_, err := orm.Limit(1).Desc("id").Get(m)
	return m, err
}

func (m *Moods) Get() (*Moods, bool) {
	has, err := orm.Get(m)
	if err != nil {
		logrus.Error(err)
		return m, false
	}
	return m, has
}

func (m *Moods) GetList(start int, num int) ([]*UserMoods, error) {
	var moods = make([]*UserMoods, 0)
	start = (start - 1) * num

	orm := orm.Select("moods.*, users.nick_name")
	orm.Join("LEFT OUTER", "users", "moods.user_id = users.id")

	orm.Desc("moods.id")
	orm.Limit(num, start)

	err := orm.Find(&moods)
	return moods, err
}

func (m *Moods) Count() (int64, error) {
	total, err := orm.Count(m)
	return total, err
}

func (m *Moods) Insert() (int64, error) {
	affected, err := orm.Insert(m)
	return affected, err
}

func (m *Moods) Update() (int64, error) {
	affected, err := orm.Id(m.Id).Update(m)
	return affected, err
}

func (m *Moods) Delete() (int64, error) {
	affected, err := orm.Id(m.Id).Delete(m)
	return affected, err
}
