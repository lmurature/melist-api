package lists

import "github.com/gin-gonic/gin"

func CreateList(c *gin.Context) {
	panic("implement me")
}

func GetListById(c *gin.Context) {
	// If list is private, we should check if caller owns the list or has access.
	panic("implement me")
}

func UpdateList(c *gin.Context) {
	// Only the owner can do this. Title, description and privacy changes.
	panic("implement me")
}

func GiveUsersAccessToList(c *gin.Context) {
	// Only the owner and those who have write access can do this. users will come in queryparam
	// ?users=123,456,678,235 separated by commas.
	panic("implement me")
}

func SearchPublicLists(c *gin.Context) {
	panic("implement me")
}

func GetMyLists(c *gin.Context) {
	panic("implement me")
}

func GetMySharedLists(c *gin.Context) {
	panic("implement me")
}
