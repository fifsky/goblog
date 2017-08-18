package models

import "time"

// table posts
type Posts struct {
	Id        uint         `xorm:"pk" json:"id"`
	CateId    uint         `xorm:"notnull" json:"cate_id"`
	Type      uint8        `xorm:"notnull" json:"type"`               // title
	UserId    uint         `xorm:"notnull" json:"user_id"`            // title
	Title     string       `xorm:"varchar(200) notnull" json:"title"` // title
	Url       string       `xorm:"varchar(100) notnull" json:"url"`   // title
	Content   string       `xorm:"longtext notnull" json:"content"`   // body
	CreatedAt time.Time    `xorm:"notnull"`
	UpdatedAt time.Time    `xorm:"notnull"`
	DeletedAt *time.Time   `xorm:"notnull"`
}

type UserPosts struct {
	Posts    `xorm:"extends"`
	Name     string
	NickName string
}

func (UserPosts) TableName() string {
	return "posts"
}

func (this *Posts) Get() (*Posts, error) {
	_, err := orm.Get(this)
	return this, err
}

func (this *Posts) GetList(start int, num int) ([]UserPosts, error) {
	var posts = make([]UserPosts, 0)
	start = (start - 1) * num
	//err := orm.Limit(num, start).Find(&posts)
	err := orm.SQL("select p.*, c.name,u.nick_name from posts p left join users u on p.user_id = u.id left join cates c on p.cate_id = c.id limit ?,?", start, num).Find(&posts)

	return posts, err
}
