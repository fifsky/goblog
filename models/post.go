package models

import (
	"time"
	"github.com/sirupsen/logrus"
)

// table posts
type Posts struct {
	Id        uint         `form:"id" xorm:"pk" json:"id"`
	CateId    uint         `form:"cate_id" xorm:"notnull" json:"cate_id"`
	Type      uint8        `form:"type" xorm:"notnull" json:"type"`
	UserId    uint         `form:"-" xorm:"notnull" json:"user_id"`
	Title     string       `form:"title" xorm:"varchar(200) notnull" json:"title"`
	Url       string       `form:"url" xorm:"varchar(100) notnull" json:"url"`
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

func (p *Posts) Get() (*Posts, bool) {
	has, err := orm.Get(p)
	if err != nil {
		logrus.Error(err)
		return p, false
	}
	return p, has
}

func (p *Posts) Prev(id uint) (*Posts, bool) {
	post := &Posts{}
	has, err := orm.Where("id < ? and type = 1", id).Desc("id").Limit(1).Get(post)
	if err != nil {
		logrus.Error(err)
		return post, false
	}
	return post, has
}

func (p *Posts) Next(id uint) (*Posts, bool) {
	post := &Posts{}
	has, err := orm.Where("id > ? and type = 1", id).Asc("id").Limit(1).Get(post)
	if err != nil {
		logrus.Error(err)
		return post, false
	}
	return post, has
}

func (p *Posts) GetList(start int, num int) ([]*UserPosts, error) {
	var posts = make([]*UserPosts, 0)
	start = (start - 1) * num
	//err := orm.Limit(num, start).Find(&posts)
	//下面这个方式拼接WHERE太麻烦
	//err := orm.SQL("select p.*, c.name,u.nick_name,c.domain from posts p left join users u on p.user_id = u.id left join cates c on p.cate_id = c.id order by p.id desc limit ?,?", start, num).Find(&posts)

	orm := orm.Select("posts.*, cates.name,users.nick_name,cates.domain")
	orm.Join("LEFT OUTER", "users", "posts.user_id = users.id")
	orm.Join("LEFT OUTER", "cates", "posts.cate_id = cates.id")

	if p.CateId > 0 {
		orm.Where("posts.cate_id = ?", p.CateId)
	}

	if p.Type > 0 {
		orm.Where("posts.type = ?", p.Type)
	}

	orm.Desc("posts.id")
	orm.Limit(num, start)

	err := orm.Find(&posts)
	return posts, err
}

func (p *Posts) Count() (int64, error) {
	total, err := orm.Count(p)
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
