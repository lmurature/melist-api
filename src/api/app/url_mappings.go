package app

import (
	auth_controller "github.com/lmurature/melist-api/src/api/controllers/auth"
	items_controller "github.com/lmurature/melist-api/src/api/controllers/items"
	lists_controller "github.com/lmurature/melist-api/src/api/controllers/lists"
	"github.com/lmurature/melist-api/src/api/controllers/ping"
	users_controller "github.com/lmurature/melist-api/src/api/controllers/users"
	"github.com/lmurature/melist-api/src/api/middlewares"
)

func mapUrls() {
	router.GET("/ping", ping.Ping)

	// TODO: https://www.npmjs.com/package/store-js para guardar access_token y refresh_token en el front
	router.POST("/api/users/auth/generate_token", auth_controller.AuthenticateUser)
	router.POST("/api/users/auth/refresh_token", auth_controller.RefreshAuthentication)

	router.GET("/api/users/me", middlewares.Authenticate(), users_controller.GetUserMe)

	router.GET("/api/items/search", middlewares.Authenticate(), items_controller.SearchItems)
	router.GET("/api/items/:item_id", middlewares.Authenticate(), items_controller.GetItem)


	// List management
	router.POST("/api/lists/create", middlewares.Authenticate(), lists_controller.CreateList)
	router.GET("/api/lists/get/:list_id", middlewares.Authenticate(), lists_controller.GetListById)
	router.GET("/api/lists/get/:list_id/shares", middlewares.Authenticate(), lists_controller.GetListShareConfigs)
	router.PUT("/api/lists/update/:list_id", middlewares.Authenticate(), lists_controller.UpdateList)
	router.PUT("/api/lists/access/:list_id", middlewares.Authenticate(), lists_controller.GiveUsersAccessToList)
	router.GET("/api/lists/search", middlewares.Authenticate(), lists_controller.SearchPublicLists)
	router.GET("/api/lists/get/all_owned", middlewares.Authenticate(), lists_controller.GetMyLists)
	router.GET("/api/lists/get/all_shared", middlewares.Authenticate(), lists_controller.GetMySharedLists)

	// List items management
	router.POST("/api/lists/:list_id/add_items/:item_id", middlewares.Authenticate(), lists_controller.AddItemsToList)
}
