package models

import (
	"time"

	"github.com/ilibs/gosql"
)

type Moods struct {
	Id        int       `form:"id" json:"id" db:"id"`
	Content   string    `form:"content" json:"content" db:"content"`
	UserId    int       `form:"user_id" json:"user_id" db:"user_id"`
	CreatedAt time.Time `form:"-" json:"created_at" db:"created_at"`
}

func (m *Moods) DbName() string {
	return "default"
}

func (m *Moods) TableName() string {
	return "moods"
}

func (m *Moods) PK() string {
	return "id"
}

func (m *Moods) AfterChange() {
	Cache.Delete("mood-first")
}

type UserMoods struct {
	Moods
	NickName string `db:"nick_name"`
}

func MoodFrist() (*UserMoods, error) {
	if v, ok := Cache.Get("mood-first"); ok {
		return v.(*UserMoods), nil
	}

	m := &UserMoods{}

	err := gosql.QueryRowx("select m.*,u.nick_name from moods m left join users u on m.user_id = u.id order by m.id desc limit 1").StructScan(m)

	if err != nil {
		return nil, err
	}

	Cache.Set("mood-first", m, 1*time.Hour)

	return m, nil
}

func MoodGetList(start int, num int) ([]*UserMoods, error) {

	var moods = make([]*UserMoods, 0)
	start = (start - 1) * num

	err := gosql.Select(&moods, "select m.*,u.nick_name from moods m left join users u on m.user_id = u.id order by m.id desc limit ?,?", start, num)

	if err != nil {
		return nil, err
	}

	return moods, err
}
