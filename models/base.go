package models

import (
	"fmt"
	"net/url"
	"github.com/fifsky/goblog/system"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"database/sql"
)

var orm *xorm.Engine

func InitDB() (*xorm.Engine, error) {
	config := system.GetConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=", config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Name)
	dsn += url.QueryEscape("Asia/Shanghai")

	var err error
	orm, err = xorm.NewEngine("mysql", dsn)
	orm.SetMaxIdleConns(200)
	orm.SetMaxOpenConns(200)
	orm.ShowSQL(true)
	return orm, err
}

func ImportDB() ([]sql.Result, error) {
	sqlpath := "./db/blog.sql"
	rst, err := orm.ImportFile(sqlpath)
	return rst, err
}
