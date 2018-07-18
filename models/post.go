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

type UserPosts struct {
	Posts
	Name     string `db:"name"`
	NickName string `db:"nick_name"`
	Domain   string `db:"domain"`
}

func PostPrev(id int) (*Posts, error) {
	m := &Posts{
		Id:   id,
		Type: 1,
	}
	err := gosql.Model(m).OrderBy("id desc").Limit(1).Get()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func PostNext(id int) (*Posts, error) {
	m := &Posts{
		Id:   id,
		Type: 1,
	}
	err := gosql.Model(m).OrderBy("id asc").Limit(1).Get()

	if err != nil {
		return nil, err
	}
	return m, nil
}

func PostArchive() ([]map[string]string, error) {
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

	return m, err
}

func PostGetList(p *Posts, start int, num int, artdate string) ([]*UserPosts, error) {
	var posts = make([]*UserPosts, 0)
	start = (start - 1) * num

	args := make([]interface{}, 0)
	where := "where 1 = 1 "

	if p.CateId > 0 {
		where += " and p.cate_id = ?"
		args = append(args, p.CateId)
	}

	if p.Type > 0 {
		where += " and p.type = ?"
		args = append(args, p.Type)
	}

	if artdate != "" {
		where += " and DATE_FORMAT(p.created_at,'%Y-%m') = ?"
		args = append(args, artdate)
	}

	args = append(args, start, num)

	rows, err := gosql.Queryx("select p.*,c.name,u.nick_name,c.domain from posts p left join users u on p.user_id = u.id left join cates c on p.cate_id = c.id "+where+" order by p.id desc limit ?,?", args...)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		m := &UserPosts{}
		err := rows.StructScan(m)
		if err != nil {
			return nil, err
		}
		posts = append(posts, m)
	}

	return posts, err
}
