package users_service

import (
	"fmt"

	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	"github.com/lmurature/melist-api/src/api/domain/users"
	users_provider "github.com/lmurature/melist-api/src/api/providers/users"
	"github.com/sirupsen/logrus"
)


type usersService struct{}

type usersServiceInterface interface {
	GetUser(userId int64, accessToken string) (*users.User, apierrors.ApiError)
}

var (
	UsersService usersServiceInterface
)

func init() {
	UsersService = &usersService{}
}

func (s *usersService) GetUser(userId int64, accessToken string) (*users.User, apierrors.ApiError) {
	user, err := users_provider.GetUserInformation(userId, accessToken)
	if err != nil {
		logrus.Error(fmt.Sprintf("error getting user %d", userId), err)
		return nil, err
	}
	return user, nil
}