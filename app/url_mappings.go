package app

import (
	auth_controller "github.com/lmurature/melist-api/controllers/auth"
	items_controller "github.com/lmurature/melist-api/controllers/items"
	users_controller "github.com/lmurature/melist-api/controllers/users"
	"github.com/lmurature/melist-api/controllers/ping"
)


func mapUrls() {
	router.GET("/ping", ping.Ping)

	router.POST("/api/users/auth/generate_token", auth_controller.AuthenticateUser)
	router.POST("/api/users/auth/refresh_token", auth_controller.RefreshAuthentication)

	router.GET("/api/users/:user_id", users_controller.GetUser)

	router.GET("/api/items/search", items_controller.SearchItems)
	router.GET("/api/items/:item_id", items_controller.GetItem)
}
