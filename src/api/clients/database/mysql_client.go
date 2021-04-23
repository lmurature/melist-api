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
	DbClient, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s/melist", config.DbUser, config.DbPass, config.DbHost))
	if err != nil {
		panic(err)
	}
	DbClient.SetConnMaxLifetime(time.Minute * 3)
	DbClient.SetMaxOpenConns(10)
	DbClient.SetMaxIdleConns(10)
}
