package app

import (
	authcontroller "github.com/lmurature/melist-api/src/api/controllers/auth"
	itemscontroller "github.com/lmurature/melist-api/src/api/controllers/items"
	ping2 "github.com/lmurature/melist-api/src/api/controllers/ping"
	userscontroller "github.com/lmurature/melist-api/src/api/controllers/users"
	"github.com/lmurature/melist-api/src/api/middlewares"
)


func mapUrls() {
	router.GET("/ping", ping2.Ping)

	router.POST("/api/users/auth/generate_token", authcontroller.AuthenticateUser)
	router.POST("/api/users/auth/refresh_token", authcontroller.RefreshAuthentication)

	router.GET("/api/users/me",  middlewares.Authenticate(), userscontroller.GetUserMe)

	router.GET("/api/items/search",  middlewares.Authenticate(), itemscontroller.SearchItems)
	router.GET("/api/items/:item_id", middlewares.Authenticate(), itemscontroller.GetItem)
}
