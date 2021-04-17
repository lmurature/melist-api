package users_controller

import (
	apierrors2 "github.com/lmurature/melist-api/src/api/domain/apierrors"
	users_provider2 "github.com/lmurature/melist-api/src/api/providers/users"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	quantityOfStringsAfterSpliting = 2
)

func GetUser(c *gin.Context) {
	paramUserId := c.Param("user_id")
	if paramUserId == "" {
		err := apierrors2.NewBadRequestApiError("'user_id' can't be empty")
		c.JSON(err.Status(), err)
		return
	}

	userId, err := strconv.ParseInt(paramUserId, 10, 64)
	if err != nil {
		err := apierrors2.NewBadRequestApiError("'user_id' must be a number")
		c.JSON(err.Status(), err)
		return
	}

	reqToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")

	if len(splitToken) != quantityOfStringsAfterSpliting {
		err := apierrors2.NewBadRequestApiError("authorization token (Bearer) is needed to access this endpoint")
		c.JSON(err.Status(), err)
		return
	}

	accessToken := splitToken[1]

	user, userErr := users_provider2.GetUserInformation(userId, accessToken)
	if userErr != nil {
		c.JSON(userErr.Status(), userErr)
		return
	}

	c.JSON(http.StatusOK, user)
}