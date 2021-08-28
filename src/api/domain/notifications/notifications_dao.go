package notifications

import (
	"github.com/lmurature/melist-api/src/api/clients/database"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	error_utils "github.com/lmurature/melist-api/src/api/utils/error"
	"github.com/sirupsen/logrus"
)

const (
	insertNotification   = "INSERT INTO list_notifications(list_id,message,timestamp,seen,permalink) VALUES(?,?,?,?,?);"
	getListNotifications = "SELECT id,list_id,message,timestamp,seen,permalink FROM list_notifications WHERE list_id=? ORDER BY timestamp DESC;"
)

var (
	NotificationsDao notificationsDaoInterface
)

type notificationsDaoInterface interface {
	SaveNotification(notification Notification) (*Notification, apierrors.ApiError)
	GetListNotifications(listId int64) ([]Notification, apierrors.ApiError)
}

type notificationsDao struct{}

func init() {
	NotificationsDao = &notificationsDao{}
}

func (n *notificationsDao) SaveNotification(notification Notification) (*Notification, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(insertNotification)
	if err != nil {
		logrus.Error("error when trying to prepare insert notification statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to insert notification", error_utils.GetDatabaseGenericError())
	}
	defer stmt.Close()

	result, saveErr := stmt.Exec(notification.ListId, notification.Message,
		notification.Timestamp, notification.Seen, notification.Permalink)
	if saveErr != nil {
		logrus.Error("error when trying to insert notification", saveErr)
		return nil, apierrors.NewInternalServerApiError("error when trying to insert notification", error_utils.GetDatabaseGenericError())
	}

	id, _ := result.LastInsertId()
	notification.Id = id

	return &notification, nil
}

func (n *notificationsDao) GetListNotifications(listId int64) ([]Notification, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(getListNotifications)
	if err != nil {
		logrus.Error("error when trying to prepare get notification statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to get notifications", error_utils.GetDatabaseGenericError())
	}
	defer stmt.Close()

	rows, err := stmt.Query(listId)
	if err != nil {
		logrus.Error("error while getting notifications", err)
		return nil, apierrors.NewInternalServerApiError("error getting notifications", error_utils.GetDatabaseGenericError())
	}
	defer rows.Close()

	result := make([]Notification, 0)
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.Id, &n.ListId, &n.Message, &n.Timestamp, &n.Seen, &n.Permalink); err != nil {
			logrus.Error("error scaning notification into notification struct", err)
			return nil, apierrors.NewInternalServerApiError("error getting notifications", error_utils.GetDatabaseGenericError())
		}
		result = append(result, n)
	}

	return result, nil
}
