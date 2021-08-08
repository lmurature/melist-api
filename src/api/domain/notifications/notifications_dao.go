package notifications

import (
	"github.com/lmurature/melist-api/src/api/clients/database"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	error_utils "github.com/lmurature/melist-api/src/api/utils/error"
	"github.com/sirupsen/logrus"
)

const (
	insertNotification = "INSERT INTO list_notifications(list_id,message,timestamp,seen,permalink) VALUES(?,?,?,?,?);"
)

var (
	NotificationsDao notificationsDaoInterface
)

type notificationsDaoInterface interface {
	SaveNotification(notification Notification) (*Notification, apierrors.ApiError)
	GetListNotifications(listId int64) ([]Notification, apierrors.ApiError)
}

type notificationsDao struct {}

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
	panic("implement me")
}


