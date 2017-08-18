package models

import "time"

// table posts
type Posts struct {
	Id        uint         `xorm:"pk" json:"id"`
	CreatedAt time.Time    `xorm:"notnull" json:"created_at"`
	UpdatedAt time.Time    `xorm:"notnull" json:"updated_at"`
	CateId    uint         `xorm:"notnull" json:"cate_id"`
	Type      uint8        `xorm:"notnull" json:"type"`               // title
	UserId    uint         `xorm:"notnull" json:"user_id"`            // title
	Title     string       `xorm:"varchar(200) notnull" json:"title"` // title
	Url       string       `xorm:"varchar(100) notnull" json:"url"`   // title
	Content   string       `xorm:"longtext notnull" json:"content"`   // body
}

func (this *Posts) Get() (*Posts, error) {
	post := &Posts{Id: this.Id}
	_, err := orm.Get(post)
	return post, err
}

func (this *Posts) GetList() ([]Posts, error) {
	var posts = make([]Posts, 0)
	err := orm.Limit(10, 3).Find(&posts)
	return posts, err
}
