package config

import (
	"os"
)

var (
	AppId int64 = 5112680121711673
	RedirectUri string
	SecretKey string

	DbDateLayout = "2006-01-02 15:04:05"
)

func init() {
	if !isDevelopment() {
		RedirectUri = "http://melist-app.herokuapp.com/auth/authorized"
	} else {
		RedirectUri = "http://localhost:3000/auth/authorized"
	}

	SecretKey = os.Getenv("SECRET_KEY")
}

func isDevelopment() bool {
	scope := os.Getenv("SCOPE")
	return scope == "development"
}
