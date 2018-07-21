package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ilibs/gosql"
)

func ImportDB() ([]sql.Result, error) {
	sqlpath := "./db/blog.sql"
	rst, err := gosql.Import(sqlpath)
	return rst, err
}