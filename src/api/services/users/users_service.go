package users_service

import (
	"fmt"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	"github.com/lmurature/melist-api/src/api/domain/lists"
	"github.com/lmurature/melist-api/src/api/domain/share"
	"github.com/lmurature/melist-api/src/api/domain/users"
	"github.com/lmurature/melist-api/src/api/providers/mail"
	users_provider "github.com/lmurature/melist-api/src/api/providers/users"
	"github.com/lmurature/melist-api/src/api/utils/date"
	"github.com/sirupsen/logrus"
)

type usersService struct{}

type usersServiceInterface interface {
	GetMeliUser(userId int64) (*users.User, apierrors.ApiError)
	GetMyUser(accessToken string) (*users.User, apierrors.ApiError)
	SaveUserToDb(u users.User, accessToken string, refreshToken string) apierrors.ApiError
	GetUserFromDb(userId int64) (*users.MelistUser, apierrors.ApiError)
	FindUserByEmail(email string) (*users.MelistUser, apierrors.ApiError)
	UpdateUserDb(u users.User, accessToken string, refreshToken string) apierrors.ApiError
	SearchUsers(query string) ([]users.MelistUser, apierrors.ApiError)
	InviteUser(email string, shareType string, listId int64, callerId int64) (map[string]interface{}, apierrors.ApiError)
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
	user := users.MelistUser{
		Id:           u.Id,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Nickname:     u.Nickname,
		Email:        u.Email,
		DateCreated:  date_utils.GetNowDateFormatted(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	_, err := users.UserDao.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *usersService) GetUserFromDb(userId int64) (*users.MelistUser, apierrors.ApiError) {
	user, err := users.UserDao.GetUser(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *usersService) FindUserByEmail(email string) (*users.MelistUser, apierrors.ApiError) {
	user, err := users.UserDao.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *usersService) UpdateUserDb(u users.User, accessToken string, refreshToken string) apierrors.ApiError {
	user := users.MelistUser{
		Id:           u.Id,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Nickname:     u.Nickname,
		Email:        u.Email,
		DateCreated:  date_utils.GetNowDateFormatted(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	_, err := users.UserDao.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *usersService) SearchUsers(query string) ([]users.MelistUser, apierrors.ApiError) {
	if query == "" {
		return nil, apierrors.NewBadRequestApiError("search users query should not be empty")
	}

	result, err := users.UserDao.SearchUsers(query)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *usersService) InviteUser(email string, shareType string, listId int64, callerId int64) (map[string]interface{}, apierrors.ApiError) {
	caller, err := s.GetUserFromDb(callerId)
	if err != nil {
		return nil, err
	}

	list, err := lists.ListDao.GetList(listId)
	if err != nil {
		return nil, err
	}

	if err := list.ValidateUpdatability(callerId); err != nil {
		return nil, err
	}

	_, err = share.ShareConfigDao.CreateEmailShareConfig(share.ShareConfig{
		ListId:    listId,
		Email:     email,
		ShareType: shareType,
	})
	if err != nil {
		return nil, err
	}

	listUrl := fmt.Sprintf("https://melist-app.herokuapp.com/lists/%d", listId)
	mail.SendMail(email, share.GetFormattedShareType(shareType),
		caller.FirstName,
		caller.LastName,
		list.Title,
		listUrl)

	return nil, nil
}

// admin of list send invitation over email to another external user
// the system creates a future collaboration row in which
// when the user that was invited logs in, the future collaboration table by email
// is read and all the share configs are persisted.
