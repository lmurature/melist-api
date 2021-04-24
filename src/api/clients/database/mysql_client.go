package database

import (
	"database/sql"
	"fmt"
	"github.com/lmurature/melist-api/src/api/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DbClient *sql.DB
)

func init() {
	var err error
	url := fmt.Sprintf("%s:%s@tcp(%s)/%s", config.DbUser, config.DbPass, config.DbHost, config.DbName)
	fmt.Println("about to connect to url ", url)
	DbClient, err = sql.Open("mysql", url)
	if err != nil {
		panic(err)
	}
	DbClient.SetConnMaxLifetime(time.Minute * 3)
	DbClient.SetMaxOpenConns(10)
	DbClient.SetMaxIdleConns(10)
}
