package database

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DbClient *sql.DB
)

func init() {
	var err error
	DbClient, err = sql.Open("mysql", "root:root@/melist")
	if err != nil {
		panic(err)
	}
	DbClient.SetConnMaxLifetime(time.Minute * 3)
	DbClient.SetMaxOpenConns(10)
	DbClient.SetMaxIdleConns(10)
}
