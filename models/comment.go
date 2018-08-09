package models

import (
	"time"
	"github.com/ilibs/gosql"
)

type Comments struct {
	Id        int       `form:"id" json:"id" db:"id"`
	PostId    int       `form:"post_id" json:"post_id" db:"post_id"`
	Pid       int       `form:"pid" json:"pid" db:"pid"`
	Name      string    `form:"name" json:"name" db:"name"`
	Content   string    `form:"content" json:"content" db:"content"`
	IP        string    `form:"-" json:"ip" db:"ip"`
	CreatedAt time.Time `form:"-" json:"created_at" db:"created_at"`
}

func (c *Comments) DbName() string {
	return "default"
}

func (c *Comments) TableName() string {
	return "links"
}

func (c *Comments) PK() string {
	return "id"
}

func PostComments(start int, num int) ([]*Comments, error) {
	var m = make([]*Comments, 0)
	start = (start - 1) * num
	err := gosql.Model(&m).OrderBy("id desc").Limit(num).Offset(start).All()
	if err != nil {
		return nil, err
	}
	return m, nil
}
