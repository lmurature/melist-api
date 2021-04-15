package app

import (
	"github.com/lmurature/melist-api/controllers/auth"
	"github.com/lmurature/melist-api/controllers/ping"
)


func mapUrls() {
	router.GET("/ping", ping.Ping)

	router.POST("/api/users/auth/generate_token", auth.AuthenticateUser)
	router.POST("/api/users/auth/refresh_token", auth.RefreshAuthentication)
}
