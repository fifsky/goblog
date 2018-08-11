package models

import (
	"time"
	"github.com/ilibs/gosql"
	"strings"
	"github.com/fifsky/goblog/helpers"
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
	return "comments"
}

func (c *Comments) PK() string {
	return "id"
}

func PostComments(postId, start, num int) ([]*Comments, error) {
	var m = make([]*Comments, 0)
	start = (start - 1) * num
	err := gosql.Model(&m).Where("post_id = ?", postId).OrderBy("id asc").Limit(num).Offset(start).All()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func PostCommentNum(postId []int) (map[int]int, error) {
	postIds := make([]string, 0)
	for _, v := range postId {
		postIds = append(postIds, helpers.ToStr(v))
	}

	rows, err := gosql.Queryx("select count(*) comment_num,post_id from comments where post_id in(" + strings.Join(postIds, ",") + ") group by post_id")

	if err != nil {
		return nil, err
	}

	m := make(map[int]int)

	for rows.Next() {
		var commentNum, postId int
		err := rows.Scan(&commentNum, &postId)
		if err != nil {
			return nil, err
		}

		m[postId] = commentNum
	}

	return m, nil
}
