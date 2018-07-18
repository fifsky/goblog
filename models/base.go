package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ilibs/gosql"
	"github.com/fifsky/goblog/system"
)

func InitDB() (error) {
	config := system.GetConfig()
	configs := make(map[string]*gosql.Config)
	configs["default"] = config.Database

	return gosql.Connect(configs)
}

func ImportDB() ([]sql.Result, error) {
	sqlpath := "./db/blog.sql"
	rst, err := gosql.Import(sqlpath)
	return rst, err
}