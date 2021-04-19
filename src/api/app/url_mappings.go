package app

import (
	auth_controller "github.com/lmurature/melist-api/src/api/controllers/auth"
	items_controller "github.com/lmurature/melist-api/src/api/controllers/items"
	"github.com/lmurature/melist-api/src/api/controllers/ping"
	users_controller "github.com/lmurature/melist-api/src/api/controllers/users"
	"github.com/lmurature/melist-api/src/api/middlewares"
)


func mapUrls() {
	router.GET("/ping", ping.Ping)

	// TODO: https://www.npmjs.com/package/store-js para guardar
	router.POST("/api/users/auth/generate_token", auth_controller.AuthenticateUser)
	router.POST("/api/users/auth/refresh_token", auth_controller.RefreshAuthentication)

	router.GET("/api/users/me",  middlewares.Authenticate(), users_controller.GetUserMe)

	router.GET("/api/items/search",  middlewares.Authenticate(), items_controller.SearchItems)
	router.GET("/api/items/:item_id", middlewares.Authenticate(), items_controller.GetItem)
}
