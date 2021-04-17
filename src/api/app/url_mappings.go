package app

import (
	auth_controller2 "github.com/lmurature/melist-api/src/api/controllers/auth"
	items_controller2 "github.com/lmurature/melist-api/src/api/controllers/items"
	ping2 "github.com/lmurature/melist-api/src/api/controllers/ping"
	users_controller2 "github.com/lmurature/melist-api/src/api/controllers/users"
)


func mapUrls() {
	router.GET("/ping", ping2.Ping)

	router.POST("/api/users/auth/generate_token", auth_controller2.AuthenticateUser)
	router.POST("/api/users/auth/refresh_token", auth_controller2.RefreshAuthentication)

	router.GET("/api/users/:user_id", users_controller2.GetUser)

	router.GET("/api/items/search", items_controller2.SearchItems)
	router.GET("/api/items/:item_id", items_controller2.GetItem)
}
