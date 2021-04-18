package users_controller

import (
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	users_service "github.com/lmurature/melist-api/src/api/services/users"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	quantityOfStringsAfterSpliting = 2
)

func GetUserMe(c *gin.Context) {
	reqToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")

	if len(splitToken) != quantityOfStringsAfterSpliting {
		err := apierrors.NewBadRequestApiError("authorization token (Bearer) is needed to access this endpoint")
		c.JSON(err.Status(), err)
		return
	}

	accessToken := splitToken[1]

	user, userErr := users_service.UsersService.GetMyUser(accessToken)
	if userErr != nil {
		c.JSON(userErr.Status(), userErr)
		return
	}

	c.JSON(http.StatusOK, user)
}