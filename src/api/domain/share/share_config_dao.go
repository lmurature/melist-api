package share

import (
	"errors"
	"github.com/lmurature/melist-api/src/api/clients/database"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	"github.com/sirupsen/logrus"
)

const (
	insertShareConfig     = "INSERT INTO share_config(user_id, list_id, `type`) VALUES (?,?,?);"
	getShareConfigsByUser = "SELECT s.user_id, s.list_id, s.`type` FROM share_config s WHERE s.user_id=?;"
	getShareConfigsByList = "SELECT s.user_id, s.list_id, s.`type` FROM share_config s WHERE s.list_id=?;"
	updateShareConfigType = "UPDATE share_config SET `type`=? WHERE (user_id=? AND list_id=?);"
	deleteShareConfig     = "DELETE FROM share_config s WHERE (s.user_id=? AND s.list_id=?);"
)

func (s *ShareConfig) Save() apierrors.ApiError {
	stmt, err := database.DbClient.Prepare(insertShareConfig)
	if err != nil {
		logrus.Error("error when trying to prepare insert share config statement", err)
		return apierrors.NewInternalServerApiError("error when trying to insert share config", errors.New("database error"))
	}
	defer stmt.Close()

	_, saveErr := stmt.Exec(s.UserId, s.ListId, s.ShareType)
	if saveErr != nil {
		return apierrors.NewInternalServerApiError("error when trying to save share config", errors.New("database error"))
	}

	return nil
}

func (s *ShareConfig) GetAllShareConfigsByUser() (ShareConfigs, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(getShareConfigsByUser)
	if err != nil {
		logrus.Error("error when trying to prepare get share config by user statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to get share configs", errors.New("database error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(s.UserId)
	if err != nil {
		logrus.Error("error while getting share config lists from user", err)
		return nil, apierrors.NewInternalServerApiError("error getting share config lists from user", errors.New("database error"))
	}
	defer rows.Close()

	result := make([]ShareConfig, 0)

	for rows.Next() {
		var conf ShareConfig
		if err := rows.Scan(&conf.UserId, &conf.ListId, &conf.ShareType); err != nil {
			logrus.Error("error when scan share config row into share config struct", err)
			return nil, apierrors.NewInternalServerApiError("error when tying to get share configs from user", errors.New("database error"))
		}
		result = append(result, conf)
	}

	if len(result) == 0 {
		return nil, apierrors.NewNotFoundApiError("no share configs found for user")
	}

	return result, nil
}

func (s *ShareConfig) GetAllSharedConfigsByList() (ShareConfigs, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(getShareConfigsByList)
	if err != nil {
		logrus.Error("error when trying to prepare get share config by list statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to get share configs", errors.New("database error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(s.ListId)
	if err != nil {
		logrus.Error("error while getting share config lists from list", err)
		return nil, apierrors.NewInternalServerApiError("error getting share config lists from list", errors.New("database error"))
	}
	defer rows.Close()

	result := make([]ShareConfig, 0)

	for rows.Next() {
		var conf ShareConfig
		if err := rows.Scan(&conf.UserId, &conf.ListId, &conf.ShareType); err != nil {
			logrus.Error("error when scan share config row into share config struct", err)
			return nil, apierrors.NewInternalServerApiError("error when tying to get share configs from list", errors.New("database error"))
		}
		result = append(result, conf)
	}

	if len(result) == 0 {
		return nil, apierrors.NewNotFoundApiError("no share configs found for list")
	}

	return result, nil
}

func (s *ShareConfig) Update() apierrors.ApiError {
	stmt, err := database.DbClient.Prepare(updateShareConfigType)
	if err != nil {
		logrus.Error("error when trying to prepare update share config statement", err)
		return apierrors.NewInternalServerApiError("error when trying to update share configs", errors.New("database error"))
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(s.ShareType, s.UserId, s.ListId)
	if updateErr != nil {
		return apierrors.NewInternalServerApiError("error when trying to update share config", errors.New("database error"))
	}

	return nil
}

func (s *ShareConfig) Delete() apierrors.ApiError {
	stmt, err := database.DbClient.Prepare(deleteShareConfig)
	if err != nil {
		logrus.Error("error when trying to prepare delete share config statement", err)
		return apierrors.NewInternalServerApiError("error when trying to delete share config", errors.New("database error"))
	}
	defer stmt.Close()

	_, deleteErr := stmt.Exec(s.UserId, s.ListId)
	if deleteErr != nil {
		return apierrors.NewInternalServerApiError("error when trying to delete share config", errors.New("database error"))
	}

	return nil
}
