package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/lmurature/melist-api/domain/apierrors"
	"github.com/lmurature/melist-api/domain/auth"
	auth_service "github.com/lmurature/melist-api/services/auth"
	"net/http"
)

func AuthenticateUser(c *gin.Context) {
	// todo: parse response del POST y devolver al cliente el auth token. Persistirlo en algun lado?
	var request auth.ClientAuthRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonErr := apierrors.NewBadRequestApiError("bad request body")
		c.JSON(http.StatusBadRequest, jsonErr)
		return
	}

	response, err := auth_service.AuthService.AuthenticateUser(request.AuthorizationCode)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func RefreshAuthentication(c *gin.Context) {
	var request auth.ClientAuthRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonErr := apierrors.NewBadRequestApiError("bad request body")
		c.JSON(http.StatusBadRequest, jsonErr)
		return
	}

	response, err := auth_service.AuthService.RefreshAuthentication(request.RefreshToken)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, response)
}
