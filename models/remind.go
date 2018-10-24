package models

import (
	"time"

	"github.com/ilibs/gosql"
)

type Reminds struct {
	Id         int       `form:"id" json:"id" db:"id"`
	Type       int       `form:"type" json:"type" db:"type"`
	At         string    `form:"at" json:"at" db:"at"`
	Content    string    `form:"content" json:"content" db:"content"`
	RemindDate time.Time `form:"remind_date" json:"remind_date" db:"remind_date" time_format:"2006-01-02 15:04:05"`
	CreatedAt  time.Time `form:"-" json:"created_at" db:"created_at"`
}

func (r *Reminds) DbName() string {
	return "default"
}

func (r *Reminds) TableName() string {
	return "reminds"
}

func (r *Reminds) PK() string {
	return "id"
}

func RemindGetList(start int, num int) ([]*Reminds, error) {
	var m = make([]*Reminds, 0)
	start = (start - 1) * num
	err := gosql.Model(&m).OrderBy("id desc").Limit(num).Offset(start).All()
	if err != nil {
		return nil, err
	}
	return m, nil
}
