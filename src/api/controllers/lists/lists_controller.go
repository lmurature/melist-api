package lists

import (
	"github.com/gin-gonic/gin"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	"github.com/lmurature/melist-api/src/api/domain/lists"
	"github.com/lmurature/melist-api/src/api/domain/share"
	lists_service "github.com/lmurature/melist-api/src/api/services/lists"
	"net/http"
	"strconv"
)

func CreateList(c *gin.Context) {
	var listDto lists.ListDto
	if err := c.ShouldBindJSON(&listDto); err != nil {
		err := apierrors.NewBadRequestApiError("invalid list json body")
		c.JSON(err.Status(), err)
		return
	}

	ownerId, _ := c.Get("user_id")
	listDto.OwnerId = ownerId.(int64)

	result, err := lists_service.ListsService.CreateList(listDto)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func GetListById(c *gin.Context) {
	// If list is private, we should check if caller owns the list or has access.
	listParam := c.Param("list_id")
	listId, err := strconv.ParseInt(listParam, 10, 64)
	if err != nil {
		br := apierrors.NewBadRequestApiError("list id must be an integer")
		c.JSON(br.Status(), br)
		return
	}

	userId, _ := c.Get("user_id")
	callerId := userId.(int64)

	result, resErr := lists_service.ListsService.GetList(listId, callerId)
	if resErr != nil {
		c.JSON(resErr.Status(), resErr)
		return
	}

	c.JSON(http.StatusOK, result)
}

func UpdateList(c *gin.Context) {
	// Only the owner can do this. Title, description and privacy changes.
	listParam := c.Param("list_id")
	listId, err := strconv.ParseInt(listParam, 10, 64)
	if err != nil {
		br := apierrors.NewBadRequestApiError("list id must be an integer")
		c.JSON(br.Status(), br)
		return
	}

	var listDto lists.ListDto
	if err := c.ShouldBindJSON(&listDto); err != nil {
		err := apierrors.NewBadRequestApiError("invalid list json body")
		c.JSON(err.Status(), err)
		return
	}
	listDto.Id = listId

	userId, _ := c.Get("user_id")
	callerId := userId.(int64)

	result, resErr := lists_service.ListsService.UpdateList(listDto, callerId)
	if resErr != nil {
		c.JSON(resErr.Status(), resErr)
		return
	}

	c.JSON(http.StatusOK, result)
}

func GiveUsersAccessToList(c *gin.Context) {
	// Only the owner and those who have write access can do this. users will come in body.
	listParam := c.Param("list_id")
	listId, err := strconv.ParseInt(listParam, 10, 64)
	if err != nil {
		br := apierrors.NewBadRequestApiError("list id must be an integer")
		c.JSON(br.Status(), br)
		return
	}

	var shareConfigs []share.ShareConfig
	if err := c.ShouldBindJSON(&shareConfigs); err != nil {
		err := apierrors.NewBadRequestApiError("invalid share config json body")
		c.JSON(err.Status(), err)
		return
	}

	userId, _ := c.Get("user_id")
	callerId := userId.(int64)

	resErr := lists_service.ListsService.GiveAccessToUsers(listId, callerId, shareConfigs)
	if resErr != nil {
		c.JSON(resErr.Status(), resErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "granted"})
}

func SearchPublicLists(c *gin.Context) {
	panic("implement me")
}

func GetMyLists(c *gin.Context) {
	userId, _ := c.Get("user_id")
	callerId := userId.(int64)

	userLists, err := lists_service.ListsService.GetMyLists(callerId)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, userLists)
}

func GetMySharedLists(c *gin.Context) {
	panic("implement me")
}
