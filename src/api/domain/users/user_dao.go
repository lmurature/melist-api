package users

import (
	"errors"
	"fmt"
	"github.com/lmurature/melist-api/src/api/clients/database"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	"github.com/sirupsen/logrus"
)

const (
	getUser     = "SELECT u.id, u.first_name, u.last_name, u.nickname, u.email, u.date_created, u.access_token, u.refresh_token FROM user u WHERE u.id=?;"
	insertUser  = "INSERT INTO user(id, first_name, last_name, email, nickname, date_created, access_token, refresh_token) VALUES(?,?,?,?,?,?,?,?);"
	findByEmail = "SELECT u.id, u.first_name, u.last_name, u.nickname, u.email, u.date_created FROM user u WHERE u.email=?;"
)

func (u *UserDto) Get() apierrors.ApiError {
	stmt, err := database.DbClient.Prepare(getUser)
	if err != nil {
		logrus.Error("error when trying to prepare get user statement", err)
		return apierrors.NewInternalServerApiError("error when trying to get user", errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(u.Id)

	if queryErr := result.Scan(&u.Id, &u.FirstName, u.LastName, u.Nickname, u.Email, u.DateCreated, u.AccessToken, u.RefreshToken); queryErr != nil {
		logrus.Error("user not found", err)
		return apierrors.NewNotFoundApiError("user not found")
	}

	return nil
}

func (u *UserDto) Save() apierrors.ApiError {
	stmt, err := database.DbClient.Prepare(insertUser)
	if err != nil {
		logrus.Error("error when trying to prepare insert user statement", err)
		return apierrors.NewInternalServerApiError("error when trying to insert user", errors.New("database error"))
	}
	defer stmt.Close()

	_, saveErr := stmt.Exec(u.Id, u.FirstName, u.LastName, u.Email, u.Nickname, u.DateCreated, u.AccessToken, u.RefreshToken)
	if saveErr != nil {
		logrus.Error("error when trying to save user", saveErr)
		return apierrors.NewInternalServerApiError("error when trying to save user", errors.New("database error"))
	}

	logrus.Info(fmt.Sprintf("successfully registered user %d", u.Id))
	return nil
}

func (u *UserDto) FindByEmail() apierrors.ApiError {
	stmt, err := database.DbClient.Prepare(findByEmail)
	if err != nil {
		logrus.Error("error when trying to prepare find user by email statement", err)
		return apierrors.NewInternalServerApiError("error when trying to find user by email", errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(u.Id)

	if queryErr := result.Scan(&u.Id, &u.FirstName, u.LastName, u.Nickname, u.Email, u.DateCreated); queryErr != nil {
		logrus.Error("email not found in user table", err)
		return apierrors.NewNotFoundApiError("email not found in user table")
	}

	return nil
}
