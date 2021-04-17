package auth_controller

import (
	"github.com/gin-gonic/gin"
	apierrors2 "github.com/lmurature/melist-api/src/api/domain/apierrors"
	auth2 "github.com/lmurature/melist-api/src/api/domain/auth"
	auth_service2 "github.com/lmurature/melist-api/src/api/services/auth"
	"net/http"
)

func AuthenticateUser(c *gin.Context) {
	var request auth2.ClientAuthRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonErr := apierrors2.NewBadRequestApiError("bad request body")
		c.JSON(http.StatusBadRequest, jsonErr)
		return
	}

	response, err := auth_service2.AuthService.AuthenticateUser(request.AuthorizationCode)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func RefreshAuthentication(c *gin.Context) {
	var request auth2.ClientAuthRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonErr := apierrors2.NewBadRequestApiError("bad request body")
		c.JSON(http.StatusBadRequest, jsonErr)
		return
	}

	response, err := auth_service2.AuthService.RefreshAuthentication(request.RefreshToken)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, response)
}
