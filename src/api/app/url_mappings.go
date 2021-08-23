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

	// Authentication management
	router.POST("/api/users/auth/generate_token", auth_controller.AuthenticateUser)
	router.POST("/api/users/auth/refresh_token", auth_controller.RefreshAuthentication)

	router.GET("/api/users/me", middlewares.Authenticate, users_controller.GetUserMe)
	router.GET("/api/users/search", middlewares.Authenticate, users_controller.SearchUsers)

	router.GET("/api/items/search", middlewares.Authenticate, items_controller.SearchItems)
	router.GET("/api/items/:item_id", middlewares.Authenticate, items_controller.GetItem)
	router.GET("/api/items/:item_id/history", middlewares.Authenticate, items_controller.GetItemHistory)

	// List management
	router.POST("/api/lists/create", middlewares.Authenticate, lists_controller.CreateList)
	router.GET("/api/lists/get/:list_id", middlewares.Authenticate, lists_controller.GetListById)
	router.GET("/api/lists/get/:list_id/shares", middlewares.Authenticate, lists_controller.GetListShareConfigs)
	router.PUT("/api/lists/update/:list_id", middlewares.Authenticate, lists_controller.UpdateList)
	router.PUT("/api/lists/access/:list_id", middlewares.Authenticate, lists_controller.GiveUsersAccessToList)
	router.DELETE("/api/lists/access/:list_id", middlewares.Authenticate, lists_controller.RevokeUserAccessToList)
	router.PUT("/api/lists/favorite/:list_id", middlewares.Authenticate, lists_controller.SetListFavorite)
	router.DELETE("/api/lists/favorite/:list_id", middlewares.Authenticate, lists_controller.UnsetListFavorite)
	router.GET("/api/lists/search", middlewares.Authenticate, lists_controller.SearchPublicLists)
	router.GET("/api/lists/get/all_owned", middlewares.Authenticate, lists_controller.GetMyLists)
	router.GET("/api/lists/get/all_shared", middlewares.Authenticate, lists_controller.GetMySharedLists)
	router.GET("/api/lists/get/favorites", middlewares.Authenticate, lists_controller.GetFavoriteLists)
	router.GET("/api/lists/get/:list_id/permissions", middlewares.Authenticate, lists_controller.GetMyPermissions)
	router.GET("/api/lists/get/:list_id/notifications", middlewares.Authenticate, lists_controller.GetListNotifications)

	// List items management
	router.POST("/api/lists/:list_id/items/:item_id", middlewares.Authenticate, lists_controller.AddItemsToList)
	router.GET("/api/lists/:list_id/items", middlewares.Authenticate, lists_controller.GetItems)
	router.DELETE("/api/lists/:list_id/items/:item_id", middlewares.Authenticate, lists_controller.DeleteItem)
	router.PUT("/api/lists/:list_id/check/:item_id", middlewares.Authenticate, lists_controller.CheckItem)
}
