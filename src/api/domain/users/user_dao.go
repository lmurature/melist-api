package users

import (
	"fmt"
	"github.com/lmurature/melist-api/src/api/clients/database"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	error_utils "github.com/lmurature/melist-api/src/api/utils/error"
	"github.com/sirupsen/logrus"
	"strings"
)

const (
	getUser     = "SELECT u.id, u.first_name, u.last_name, u.nickname, u.email, u.date_created, u.access_token, u.refresh_token FROM user u WHERE u.id=?;"
	insertUser  = "INSERT INTO user(id, first_name, last_name, email, nickname, date_created, access_token, refresh_token) VALUES(?,?,?,?,?,?,?,?);"
	findByEmail = "SELECT u.id, u.first_name, u.last_name, u.nickname, u.email, u.date_created FROM user u WHERE u.email=?;"
	updateUser  = "UPDATE user SET first_name=?, last_name=?, email=?, nickname=?, access_token=?, refresh_token=? WHERE id=?"
)

var (
	UserDao userDaoInterface
)

type userDaoInterface interface {
	GetUser(userId int64) (*MelistUser, apierrors.ApiError)
	CreateUser(user MelistUser) (*MelistUser, apierrors.ApiError)
	GetByEmail(email string) (*MelistUser, apierrors.ApiError)
	UpdateUser(user MelistUser) (*MelistUser, apierrors.ApiError)
}

type userDao struct{}

func init() {
	UserDao = &userDao{}
}

func (dao *userDao) GetUser(userId int64) (*MelistUser, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(getUser)
	if err != nil {
		logrus.Error("error when trying to prepare get user statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to get user", error_utils.GetDatabaseGenericError())
	}
	defer stmt.Close()

	result := stmt.QueryRow(userId)

	var u MelistUser
	if queryErr := result.Scan(&u.Id, &u.FirstName, u.LastName, u.Nickname, u.Email, u.DateCreated, u.AccessToken, u.RefreshToken); queryErr != nil {
		logrus.Error("user not found", err)
		return nil, apierrors.NewNotFoundApiError("user not found")
	}

	return &u, nil
}

func (dao *userDao) CreateUser(user MelistUser) (*MelistUser, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(insertUser)
	if err != nil {
		logrus.Error("error when trying to prepare insert user statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to insert user", error_utils.GetDatabaseGenericError())
	}
	defer stmt.Close()

	_, saveErr := stmt.Exec(user.Id, user.FirstName, user.LastName, user.Email, user.Nickname, user.DateCreated, user.AccessToken, user.RefreshToken)
	if saveErr != nil {
		if !strings.Contains(saveErr.Error(), "Duplicate entry") {
			logrus.Error("error when trying to save user", saveErr)
		}
		return nil, apierrors.NewInternalServerApiError("error when trying to save user", error_utils.GetDatabaseGenericError())
	}

	logrus.Info(fmt.Sprintf("successfully registered user %d", user.Id))
	return &user, nil
}

func (dao *userDao) GetByEmail(email string) (*MelistUser, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(findByEmail)
	if err != nil {
		logrus.Error("error when trying to prepare find user by email statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to find user by email", error_utils.GetDatabaseGenericError())
	}
	defer stmt.Close()

	result := stmt.QueryRow(email)

	var user MelistUser
	if queryErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Nickname, &user.Email, &user.DateCreated); queryErr != nil {
		logrus.Error("email not found in user table", err)
		return nil, apierrors.NewNotFoundApiError("email not found in user table")
	}

	return &user, nil
}

func (dao *userDao) UpdateUser(user MelistUser) (*MelistUser, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(updateUser)
	if err != nil {
		logrus.Error("error when trying to prepare update user statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to update user", error_utils.GetDatabaseGenericError())
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Nickname, user.AccessToken, user.RefreshToken, user.Id)
	if updateErr != nil {
		logrus.Error("error when trying to update user", updateErr)
		return nil, apierrors.NewInternalServerApiError("error when trying to update user", error_utils.GetDatabaseGenericError())
	}

	logrus.Info(fmt.Sprintf("successfully updated user %d", user.Id))
	return &user, nil
}
