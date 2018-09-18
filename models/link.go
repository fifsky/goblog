package models

import (
	"time"

	"github.com/ilibs/gosql"
	"github.com/ilibs/logger"
)

type Links struct {
	Id        int       `form:"id" json:"id" db:"id"`
	Name      string    `form:"name" json:"name" db:"name"`
	Url       string    `form:"url" json:"url" db:"url"`
	Desc      string    `form:"desc" json:"desc" db:"desc"`
	CreatedAt time.Time `form:"-" json:"created_at" db:"created_at"`
}

func (l *Links) DbName() string {
	return "default"
}

func (l *Links) TableName() string {
	return "links"
}

func (l *Links) PK() string {
	return "id"
}

func (p *Links) AfterChange()  {
	Cache.Delete("all-links")
}

func LinkGetList(start int, num int) ([]*Links, error) {
	var m = make([]*Links, 0)
	start = (start - 1) * num
	err := gosql.Model(&m).OrderBy("id desc").Limit(num).Offset(start).All()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func GetAllLinks() []*Links {
	if v, ok := Cache.Get("all-links"); ok {
		return v.([]*Links)
	}
	links := make([]*Links, 0)
	err := gosql.Model(&links).All()
	if err != nil {
		logger.Error(err)
		return nil
	}

	Cache.Set("all-links", links, 1*time.Hour)
	return links
}
