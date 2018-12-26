package models

import (
	"time"

	"github.com/ilibs/gosql"
)

type Posts struct {
	Id        int       `form:"id" json:"id" db:"id"`
	CateId    int       `form:"cate_id" json:"cate_id" db:"cate_id"`
	Type      int       `form:"type" json:"type" db:"type"`
	UserId    int       `form:"user_id" json:"user_id" db:"user_id"`
	Title     string    `form:"title" json:"title" db:"title"`
	Url       string    `form:"url" json:"url" db:"url"`
	Content   string    `form:"content" json:"content" db:"content"`
	ViewNum   string    `form:"view_num" json:"view_num" db:"view_num"`
	CreatedAt time.Time `form:"-" json:"created_at" db:"created_at"`
	UpdatedAt time.Time `form:"-" json:"updated_at" db:"updated_at"`
}

func (p *Posts) DbName() string {
	return "default"
}

func (p *Posts) TableName() string {
	return "posts"
}

func (p *Posts) PK() string {
	return "id"
}

func (p *Posts) AfterChange() {
	Cache.Delete("post-archive")
	Cache.Delete("all-cates")
}

type UserPosts struct {
	Posts
	Cate       *Cates `db:"-" relation:"cate_id,id"`
	User       *Users `db:"-" relation:"user_id,id"`
	CommentNum int    `db:"-"`
}

func PostPrev(id int) (*Posts, error) {
	m := &Posts{
		Type: 1,
	}
	err := gosql.Model(m).Where("id < ?", id).OrderBy("id desc").Limit(1).Get()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func PostNext(id int) (*Posts, error) {
	m := &Posts{
		Type: 1,
	}
	err := gosql.Model(m).Where("id > ?", id).OrderBy("id asc").Limit(1).Get()

	if err != nil {
		return nil, err
	}
	return m, nil
}

func PostArchive() ([]map[string]string, error) {

	if v, ok := Cache.Get("post-archive"); ok {
		return v.([]map[string]string), nil
	}

	m := make([]map[string]string, 0)
	result, err := gosql.Queryx("select ym,count(ym) total from (select DATE_FORMAT(created_at,'%Y/%m') as ym from posts where type = 1) s group by ym order by ym desc")

	if err != nil {
		return nil, err
	}

	for result.Next() {
		var ym, total string
		result.Scan(&ym, &total)
		m = append(m, map[string]string{
			"ym":    ym,
			"total": total,
		})
	}

	Cache.Set("post-archive", m, 1*time.Hour)

	return m, err
}

func PostGetList(p *Posts, start int, num int, artdate, keyword string) ([]*UserPosts, error) {
	var posts = make([]*UserPosts, 0)
	start = (start - 1) * num

	args := make([]interface{}, 0)
	where := "1 = 1 "

	if p.CateId > 0 {
		where += " and cate_id = ?"
		args = append(args, p.CateId)
	}

	if p.Type > 0 {
		where += " and type = ?"
		args = append(args, p.Type)
	}

	if artdate != "" {
		where += " and DATE_FORMAT(created_at,'%Y-%m') = ?"
		args = append(args, artdate)
	}

	if keyword != "" {
		where += " and title like ?"
		args = append(args, "%"+keyword+"%")
	}

	err := gosql.Model(&posts).Where(where, args...).Limit(num).Offset(start).OrderBy("id desc").All()

	if err != nil {
		return nil, err
	}

	postIds := make([]int, 0)

	for _,v := range posts {
		postIds = append(postIds, v.Id)
	}

	cm, err := PostCommentNum(postIds)

	if err != nil {
		return nil, err
	}

	for _, v := range posts {
		if c, ok := cm[v.Id]; ok {
			v.CommentNum = c
		}
	}

	return posts, err
}
