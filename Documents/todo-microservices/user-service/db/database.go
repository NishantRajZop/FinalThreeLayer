package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB(filepath string) {
	var err error
	DB, err = sql.Open("mysql", filepath)
	if err != nil {
		panic(err)
	}

}
