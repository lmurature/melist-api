package users_service

import (
	"fmt"
	"github.com/lmurature/melist-api/src/api/config"
	"time"

	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	"github.com/lmurature/melist-api/src/api/domain/users"
	users_provider "github.com/lmurature/melist-api/src/api/providers/users"
	"github.com/sirupsen/logrus"
)

type usersService struct{}

type usersServiceInterface interface {
	GetMeliUser(userId int64) (*users.User, apierrors.ApiError)
	GetMyUser(accessToken string) (*users.User, apierrors.ApiError)
	SaveUserToDb(u users.User, accessToken string, refreshToken string) apierrors.ApiError
	GetUserFromDb(userId int64) (*users.UserDto, apierrors.ApiError)
	FindUserByEmail(email string) (*users.UserDto, apierrors.ApiError)
	UpdateUserDb(u users.User, accessToken string, refreshToken string) apierrors.ApiError
}

var (
	UsersService usersServiceInterface
)

func init() {
	UsersService = &usersService{}
}

func (s *usersService) GetMeliUser(userId int64) (*users.User, apierrors.ApiError) {
	user, err := users_provider.GetUserInformation(userId)
	if err != nil {
		logrus.Error(fmt.Sprintf("error getting user %d", userId), err)
		return nil, err
	}
	return user, nil
}

func (s *usersService) GetMyUser(accessToken string) (*users.User, apierrors.ApiError) {
	user, err := users_provider.GetUserInformationMe(accessToken)
	if err != nil {
		logrus.Error("error getting my user information", err)
		return nil, err
	}
	return user, nil
}

func (s *usersService) SaveUserToDb(u users.User, accessToken string, refreshToken string) apierrors.ApiError {
	userDto := users.UserDto{
		Id:           u.Id,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Nickname:     u.Nickname,
		Email:        u.Email,
		DateCreated:  time.Now().UTC().Format(config.DbDateLayout),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return userDto.Save()
}

func (s *usersService) GetUserFromDb(userId int64) (*users.UserDto, apierrors.ApiError) {
	userDto := users.UserDto{Id: userId}
	if err := userDto.Get(); err != nil {
		return nil, err
	}
	return &userDto, nil
}

func (s *usersService) FindUserByEmail(email string) (*users.UserDto, apierrors.ApiError) {
	userDto := users.UserDto{Email: email}
	if err := userDto.FindByEmail(); err != nil {
		return nil, err
	}
	return &userDto, nil
}

func (s *usersService) UpdateUserDb(u users.User, accessToken string, refreshToken string) apierrors.ApiError {
	userDto := users.UserDto{
		Id:           u.Id,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Nickname:     u.Nickname,
		Email:        u.Email,
		DateCreated:  time.Now().UTC().Format(config.DbDateLayout),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return userDto.Update()
}
