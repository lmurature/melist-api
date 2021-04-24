package config

import (
	"fmt"
	"os"
)

var (
	AppId       int64 = 5112680121711673
	RedirectUri string
	SecretKey   string

	DbUser string
	DbPass string
	DbHost string
	DbName string

	ApiPort string

	DbDateLayout = "2006-01-02 15:04:05"
)

func init() {
	if !isDevelopment() {
		RedirectUri = "http://melist-app.herokuapp.com/auth/authorized"
		DbUser = os.Getenv("DB_USER")
		DbPass = os.Getenv("DB_PASS")
		DbHost = os.Getenv("DB_HOST")
		DbName = os.Getenv("DB_NAME")
		ApiPort = fmt.Sprintf(":%s", os.Getenv("PORT"))
	} else {
		RedirectUri = "http://localhost:3000/auth/authorized"
		DbUser = "root"
		DbPass = "root"
		DbName = "melist"
		ApiPort = ":8080"
	}

	SecretKey = os.Getenv("SECRET_KEY")
}

func isDevelopment() bool {
	scope := os.Getenv("SCOPE")
	return scope == "development"
}
