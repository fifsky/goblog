package models

import "time"

type Moods struct {
	Id        uint        `xorm:"pk"`
	Content   string      `xorm:"varchar(255) notnull"`
	UserId    uint        `xorm:"notnull"`
	CreatedAt time.Time   `xorm:"notnull"`
}

func (m *Moods) Frist() (*Moods, error) {
	_, err := orm.Limit(1).Desc("id").Get(m)
	return m, err
}