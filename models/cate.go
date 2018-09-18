package models

import (
	"time"

	"github.com/ilibs/gosql"
	"github.com/ilibs/logger"
)

type Cates struct {
	Id        int       `form:"id" json:"id" db:"id"`
	Name      string    `form:"name" json:"name" db:"name"`
	Desc      string    `form:"desc" json:"desc" db:"desc"`
	Domain    string    `form:"domain" json:"domain" db:"domain"`
	CreatedAt time.Time `form:"-" json:"created_at" db:"created_at"`
	UpdatedAt time.Time `form:"-" json:"updated_at" db:"updated_at"`
}

func (c *Cates) DbName() string {
	return "default"
}

func (c *Cates) TableName() string {
	return "cates"
}

func (c *Cates) PK() string {
	return "id"
}

func (c *Cates) AfterChange()  {
	Cache.Delete("all-cates")
}

type CateArtivleCount struct {
	Cates
	Num int `db:"num"`
}

func CateGetList(start int, num int) ([]*Cates, error) {
	var m = make([]*Cates, 0)
	start = (start - 1) * num
	err := gosql.Model(&m).OrderBy("id desc").Limit(num).Offset(start).All()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func CateArtivleCountGetList(start int, num int) ([]*CateArtivleCount, error) {
	var m = make([]*CateArtivleCount, 0)
	start = (start - 1) * num
	err := gosql.Select(&m, "select c.*,ifnull(a.num,0) num from cates c left join (select count(*) num,cate_id from posts group by cate_id) a on c.id = a.cate_id order by c.id desc limit ?,?", start, num)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func GetAllCates() []*Cates {
	if v, ok := Cache.Get("all-cates"); ok {
		return v.([]*Cates)
	}
	cates := make([]*Cates, 0)
	err := gosql.Model(&cates).All()
	if err != nil {
		logger.Error(err)
		return nil
	}

	Cache.Set("all-cates", cates, 1*time.Hour)
	return cates
}
