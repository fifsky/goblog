package models

import (
	"fmt"
	"time"
	"net/url"
	"github.com/fifsky/goblog/system"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

// I don't need soft delete,so I use customized BaseModel instead gorm.Model
type BaseModel struct {
	Id        uint `xorm:"pk"`
	CreatedAt time.Time `xorm:"notnull"`
	UpdatedAt time.Time `xorm:"notnull"`
}

var engine *xorm.Engine

func InitDB() (*xorm.Engine, error) {
	config := system.GetConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=", config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Name)
	dsn += url.QueryEscape("Asia/Shanghai")

	var err error
	engine, err = xorm.NewEngine("mysql", dsn)
	engine.SetMaxIdleConns(20)
	engine.SetMaxOpenConns(20)
	engine.ShowSQL(true)

	return engine, err
}
