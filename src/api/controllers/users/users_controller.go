package users_controller

import (
	"github.com/gin-gonic/gin"
	users_service "github.com/lmurature/melist-api/src/api/services/users"
	"net/http"
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
