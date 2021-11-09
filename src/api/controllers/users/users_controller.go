package users_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	users_service "github.com/lmurature/melist-api/src/api/services/users"
	"net/http"
	"strconv"
)


func GetUserMe(c *gin.Context) {
	token, _ := c.Get("token")
	user, userErr := users_service.UsersService.GetMyUser(token.(string))
	if userErr != nil {
		c.JSON(userErr.Status(), userErr)
		return
	}

	c.JSON(http.StatusOK, user)
}

func SearchUsers(c *gin.Context) {
	result, err := users_service.UsersService.SearchUsers(c.Query("q"))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func InviteUser(c *gin.Context) {
	email := c.Query("email")
	shareType := c.Query("share_type")
	listIdQuery := c.Query("list_id")
	listId, err := strconv.ParseInt(listIdQuery, 10, 64)
	if err != nil {
		br := apierrors.NewBadRequestApiError("list id must be an integer")
		c.JSON(br.Status(), br)
		return
	}
	callerId, _ := c.Get("user_id")

	result, resErr := users_service.UsersService.InviteUser(email, shareType, listId, callerId.(int64))
	if resErr != nil {
		c.JSON(resErr.Status(), resErr)
		return
	}

	c.JSON(http.StatusOK, result)
}
