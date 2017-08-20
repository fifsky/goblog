package models

import (
	"time"
)

// table posts
type Posts struct {
	Id        uint         `form:"id" xorm:"pk" json:"id"`
	CateId    uint         `form:"cate_id" xorm:"notnull" json:"cate_id"`
	Type      uint8        `form:"type" xorm:"notnull" json:"type"`
	UserId    uint         `form:"-" xorm:"notnull" json:"user_id"`
	Title     string       `form:"title" xorm:"varchar(200) notnull" json:"title"`
	Url       string       `form:"-" xorm:"varchar(100) notnull" json:"url"`
	Content   string       `form:"content" xorm:"longtext notnull" json:"content"`
	CreatedAt time.Time    `form:"-" xorm:"created notnull"`
	UpdatedAt time.Time    `form:"-" xorm:"updated notnull"`
}

type UserPosts struct {
	Posts    `xorm:"extends"`
	Name     string
	NickName string
	Domain   string
}

func (UserPosts) TableName() string {
	return "posts"
}

func (p *Posts) Get() (*Posts, error) {
	_, err := orm.Get(p)
	return p, err
}

func (p *Posts) GetList(start int, num int) ([]*UserPosts, error) {
	var posts = make([]*UserPosts, 0)
	start = (start - 1) * num
	//err := orm.Limit(num, start).Find(&posts)
	err := orm.SQL("select p.*, c.name,u.nick_name,c.domain from posts p left join users u on p.user_id = u.id left join cates c on p.cate_id = c.id order by p.id desc limit ?,?", start, num).Find(&posts)

	return posts, err
}

func (p *Posts) Count() (int64, error) {
	total, err := orm.Where("id >?", 1).Count(p)
	return total, err
}

func (p *Posts) Insert() (int64, error) {
	affected, err := orm.Insert(p)
	return affected, err
}

func (p *Posts) Update() (int64, error) {
	affected, err := orm.Id(p.Id).Update(p)
	return affected, err
}

func (p *Posts) Delete() (int64, error) {
	affected, err := orm.Id(p.Id).Delete(p)
	return affected, err
}
