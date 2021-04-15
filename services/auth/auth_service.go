package auth_service

import (
	"fmt"
	"github.com/lmurature/melist-api/domain/apierrors"
	"github.com/lmurature/melist-api/domain/auth"
	auth_provider "github.com/lmurature/melist-api/providers/auth"
	"github.com/sirupsen/logrus"
)

type authService struct {}

type authServiceInterface interface{
	AuthenticateUser(code string) (*auth.MeliAuthResponse, apierrors.ApiError)
	RefreshAuthentication(refreshToken string) (*auth.MeliAuthResponse, apierrors.ApiError)
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
	// TODO: Analyze if persisting this data would be useful.
	// TODO: Register user into database
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


