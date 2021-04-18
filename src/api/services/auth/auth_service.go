package auth_service

import (
	"fmt"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	"github.com/lmurature/melist-api/src/api/domain/auth"
	auth_provider "github.com/lmurature/melist-api/src/api/providers/auth"
	users_provider "github.com/lmurature/melist-api/src/api/providers/users"
	users_service "github.com/lmurature/melist-api/src/api/services/users"
	"github.com/sirupsen/logrus"
)

type authService struct{}

type authServiceInterface interface {
	AuthenticateUser(code string) (*auth.MeliAuthResponse, apierrors.ApiError)
	RefreshAuthentication(refreshToken string) (*auth.MeliAuthResponse, apierrors.ApiError)
	ValidateAccessToken(accessToken string) apierrors.ApiError
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
		// save user gives error 'cause it already exist.
		if err := users_service.UsersService.UpdateUserDb(*authenticatedUser, result.AccessToken, result.RefreshToken); err != nil {
			return nil, err
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

func (s *authService) ValidateAccessToken(accessToken string) apierrors.ApiError {
	_, err := users_provider.GetUserInformationMe(accessToken)
	return err
}
