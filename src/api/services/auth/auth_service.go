package auth_service

import (
	"fmt"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	"github.com/lmurature/melist-api/src/api/domain/auth"
	"github.com/lmurature/melist-api/src/api/domain/share"
	"github.com/lmurature/melist-api/src/api/domain/users"
	auth_provider "github.com/lmurature/melist-api/src/api/providers/auth"
	users_service "github.com/lmurature/melist-api/src/api/services/users"
	"github.com/sirupsen/logrus"
)

type authService struct{}

type authServiceInterface interface {
	AuthenticateUser(code string) (*auth.MeliAuthResponse, apierrors.ApiError)
	RefreshAuthentication(refreshToken string) (*auth.MeliAuthResponse, apierrors.ApiError)
	ValidateAccessToken(accessToken string) (*users.User, apierrors.ApiError)
}

var (
	AuthService authServiceInterface
)

func init() {
	AuthService = &authService{}
}

func (s *authService) AuthenticateUser(code string) (*auth.MeliAuthResponse, apierrors.ApiError) {
	result, err := auth_provider.CreateUserAccessToken(code)
	if err != nil {
		logrus.Error("error getting user authentication access token", err)
		return nil, err
	}

	logrus.Info(fmt.Sprintf("successfully authenticated user %d", result.UserId))

	authenticatedUser, err := users_service.UsersService.GetMyUser(result.AccessToken)
	if err != nil {
		logrus.Error("error while retrieving user information upon login")
		return nil, err
	}

	if err := users_service.UsersService.SaveUserToDb(*authenticatedUser, result.AccessToken, result.RefreshToken); err != nil {
		// save user gives error 'cause it already exist. This should occur only if client loses refresh token.

		if err := users_service.UsersService.UpdateUserDb(*authenticatedUser, result.AccessToken, result.RefreshToken); err != nil {
			return nil, err
		}

	} else {
		// parse all user email requests to collaborate to list
		userFutureConfigs, err := share.ShareConfigDao.GetAllFutureListCollaborationByEmail(authenticatedUser.Email)
		if err != nil {
			return nil, err
		}

		for _, fc := range userFutureConfigs {
			_, err := share.ShareConfigDao.CreateShareConfig(share.ShareConfig{
				ListId:    fc.ListId,
				UserId:    authenticatedUser.Id,
				ShareType: fc.ShareType,
			})
			if err == nil {
				logrus.Info(fmt.Sprintf("successfully added future collaboration as proper collaboration to user %s. Deleting future collaboration register", authenticatedUser.Email))
				err := share.ShareConfigDao.DeleteFutureCollaborationConfig(authenticatedUser.Email, fc.ListId)
				if err == nil {
					logrus.Info(fmt.Sprintf("successfully deleted future collaboration %s.", authenticatedUser.Email))
				}
			}
		}
	}

	return result, nil
}

func (s *authService) RefreshAuthentication(refreshToken string) (*auth.MeliAuthResponse, apierrors.ApiError) {
	result, err := auth_provider.RefreshAccessToken(refreshToken)
	if err != nil {
		logrus.Error("error refreshing user authentication access token", err)
		return nil, err
	}
	logrus.Info(fmt.Sprintf("successfully refreshed token for user %d", result.UserId))
	return result, nil
}

func (s *authService) ValidateAccessToken(accessToken string) (*users.User, apierrors.ApiError) {
	return users_service.UsersService.GetMyUser(accessToken)
}
