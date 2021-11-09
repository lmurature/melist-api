package share

import (
	"errors"
	"github.com/lmurature/melist-api/src/api/clients/database"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	"github.com/lmurature/melist-api/src/api/domain/users"
	"github.com/sirupsen/logrus"
)

const (
	insertShareConfig             = "INSERT INTO share_config(user_id, list_id, `type`) VALUES (?,?,?);"
	insertEmailShareConfig        = "INSERT INTO future_colaborator(user_email, list_id, share_type) VALUES (?,?,?);"
	getShareConfigsByUser         = "SELECT s.user_id, s.list_id, s.`type` FROM share_config s WHERE s.user_id=?;"
	getFutureCollaborationsByUser = "SELECT s.user_email, s.list_id, s.share_type FROM future_colaborator s WHERE s.user_email=?;"
	getShareConfigsByList         = "SELECT s.user_id, s.list_id, s.`type`, u.first_name, u.last_name, u.email, u.nickname FROM share_config s INNER JOIN user u ON s.user_id=u.id WHERE s.list_id=?;"
	updateShareConfigType         = "UPDATE share_config SET `type`=? WHERE (user_id=? AND list_id=?);"
	deleteShareConfig             = "DELETE FROM share_config s WHERE (s.user_id=? AND s.list_id=?);"
)

var (
	ShareConfigDao shareConfigDaoInterface
)

type shareConfigDaoInterface interface {
	CreateShareConfig(conf ShareConfig) (*ShareConfig, apierrors.ApiError)
	GetAllShareConfigsByUser(userId int64) (ShareConfigs, apierrors.ApiError)
	GetAllShareConfigsByList(listId int64) (ShareConfigs, apierrors.ApiError)
	UpdateShareConfig(conf ShareConfig) (*ShareConfig, apierrors.ApiError)
	DeleteShareConfig(userId int64, listId int64) apierrors.ApiError
	CreateEmailShareConfig(conf ShareConfig) (*ShareConfig, apierrors.ApiError)
	GetAllFutureListCollaborationByEmail(email string) (ShareConfigs, apierrors.ApiError)
}

type shareConfigDao struct{}

func init() {
	ShareConfigDao = &shareConfigDao{}
}

func (dao *shareConfigDao) CreateShareConfig(conf ShareConfig) (*ShareConfig, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(insertShareConfig)
	if err != nil {
		logrus.Error("error when trying to prepare insert share config statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to insert share config", errors.New("database error"))
	}
	defer stmt.Close()

	_, saveErr := stmt.Exec(conf.UserId, conf.ListId, conf.ShareType)
	if saveErr != nil {
		logrus.Error(saveErr)
		return nil, apierrors.NewInternalServerApiError("error when trying to save share config", errors.New("database error"))
	}

	return &conf, nil
}

func (dao *shareConfigDao) CreateEmailShareConfig(conf ShareConfig) (*ShareConfig, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(insertEmailShareConfig)
	if err != nil {
		logrus.Error("error when trying to prepare insert share config statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to insert share config", errors.New("database error"))
	}
	defer stmt.Close()

	_, saveErr := stmt.Exec(conf.Email, conf.ListId, conf.ShareType)
	if saveErr != nil {
		logrus.Error(saveErr)
		return nil, apierrors.NewInternalServerApiError("error when trying to save share config", errors.New("database error"))
	}

	return &conf, nil
}

func (dao *shareConfigDao) GetAllShareConfigsByUser(userId int64) (ShareConfigs, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(getShareConfigsByUser)
	if err != nil {
		logrus.Error("error when trying to prepare get share config by user statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to get share configs", errors.New("database error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId)
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

	return result, nil
}

func (dao *shareConfigDao) GetAllFutureListCollaborationByEmail(email string) (ShareConfigs, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(getFutureCollaborationsByUser)
	if err != nil {
		logrus.Error("error when trying to prepare get email share config by email statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to get email configs", errors.New("database error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(email)
	if err != nil {
		logrus.Error("error while getting email config lists from email", err)
		return nil, apierrors.NewInternalServerApiError("error getting email config lists from email", errors.New("database error"))
	}
	defer rows.Close()

	result := make([]ShareConfig, 0)

	for rows.Next() {
		var conf ShareConfig
		if err := rows.Scan(&conf.Email, &conf.ListId, &conf.ShareType); err != nil {
			logrus.Error("error when scan email config row into share config struct", err)
			return nil, apierrors.NewInternalServerApiError("error when tying to get email configs from email", errors.New("database error"))
		}
		result = append(result, conf)
	}

	return result, nil
}

func (dao *shareConfigDao) GetAllShareConfigsByList(listId int64) (ShareConfigs, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(getShareConfigsByList)
	if err != nil {
		logrus.Error("error when trying to prepare get share config by list statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to get share configs", errors.New("database error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(listId)
	if err != nil {
		logrus.Error("error while getting share config lists from list", err)
		return nil, apierrors.NewInternalServerApiError("error getting share config lists from list", errors.New("database error"))
	}
	defer rows.Close()

	result := make([]ShareConfig, 0)

	for rows.Next() {
		var conf ShareConfig
		var userData users.MelistUser
		if err := rows.Scan(&conf.UserId, &conf.ListId, &conf.ShareType, &userData.FirstName, &userData.LastName, &userData.Email, &userData.Nickname); err != nil {
			logrus.Error("error when scan share config row into share config struct", err)
			return nil, apierrors.NewInternalServerApiError("error when tying to get share configs from list", errors.New("database error"))
		}
		conf.UserData = &userData
		result = append(result, conf)
	}

	if len(result) == 0 {
		return nil, apierrors.NewNotFoundApiError("no share configs found for list")
	}

	return result, nil
}

func (dao *shareConfigDao) UpdateShareConfig(conf ShareConfig) (*ShareConfig, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(updateShareConfigType)
	if err != nil {
		logrus.Error("error when trying to prepare update share config statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to update share configs", errors.New("database error"))
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(conf.ShareType, conf.UserId, conf.ListId)
	if updateErr != nil {
		return nil, apierrors.NewInternalServerApiError("error when trying to update share config", errors.New("database error"))
	}

	return &conf, nil
}

func (dao *shareConfigDao) DeleteShareConfig(userId int64, listId int64) apierrors.ApiError {
	stmt, err := database.DbClient.Prepare(deleteShareConfig)
	if err != nil {
		logrus.Error("error when trying to prepare delete share config statement", err)
		return apierrors.NewInternalServerApiError("error when trying to delete share config", errors.New("database error"))
	}
	defer stmt.Close()

	_, deleteErr := stmt.Exec(userId, listId)
	if deleteErr != nil {
		return apierrors.NewInternalServerApiError("error when trying to delete share config", errors.New("database error"))
	}

	return nil
}
