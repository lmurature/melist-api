package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	auth_service "github.com/lmurature/melist-api/src/api/services/auth"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqToken := c.Request.Header.Get("Authorization")

		if reqToken == "" {
			apierror := apierrors.NewForbiddenApiError("Authorization token not provided")
			logrus.Error(apierror.Error(), apierror)
			c.JSON(http.StatusForbidden, apierror)
			c.Abort()
			return
		}

		splitToken := strings.Split(reqToken, "Bearer ")

		if len(splitToken) != 2 {
			err := apierrors.NewBadRequestApiError("authorization token (Bearer) is needed to access this endpoint")
			c.JSON(err.Status(), err)
			return
		}

		token := splitToken[1]

		user, err := auth_service.AuthService.ValidateAccessToken(token)

		if err != nil {
			apierror := apierrors.NewForbiddenApiError("access token not found")
			logrus.Error(apierror.Error(), apierror)
			c.JSON(http.StatusForbidden, apierror)
			c.Abort()
			return
		}

		c.Set("token", token)
		c.Set("user_id", user.Id)
		c.Next()
	}
}
