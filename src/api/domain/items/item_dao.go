package items

import (
	"fmt"
	"github.com/lmurature/melist-api/src/api/clients/database"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	error_utils "github.com/lmurature/melist-api/src/api/utils/error"
	"github.com/sirupsen/logrus"
)

const (
	insertItem = "INSERT INTO item(item_id) VALUES(?);"
	getAllItems = "SELECT i.item_id FROM item i;"
)

var (
	ItemDao itemDaoInterface
)

type itemDaoInterface interface {
	InsertItem(itemId string) apierrors.ApiError
	GetAllItems() ([]string, apierrors.ApiError)
}

type itemDao struct{}

func init() {
	ItemDao = &itemDao{}
}

func (dao *itemDao) InsertItem(itemId string) apierrors.ApiError {
	stmt, err := database.DbClient.Prepare(insertItem)
	if err != nil {
		logrus.Error("error when trying to prepare insert item statement", err)
		return apierrors.NewInternalServerApiError("error when trying to insert item", error_utils.GetDatabaseGenericError())
	}
	defer stmt.Close()

	_, saveErr := stmt.Exec(itemId)
	if saveErr != nil {
		return apierrors.NewInternalServerApiError("error when trying to save item to items table", error_utils.GetDatabaseGenericError())
	}

	logrus.Info(fmt.Sprintf("successfully added %s to item table", itemId))
	return nil
}

func (dao *itemDao) GetAllItems() ([]string, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(getAllItems)
	if err != nil {
		logrus.Error("error when trying to prepare get all items statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to get all items", error_utils.GetDatabaseGenericError())
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		logrus.Error("error while getting all items", err)
		return nil, apierrors.NewInternalServerApiError("error getting all items", error_utils.GetDatabaseGenericError())
	}
	defer rows.Close()

	result := make([]string, 0)

	for rows.Next() {
		var itemId string
		if err := rows.Scan(&itemId); err != nil {
			logrus.Error("error when scan list row into itemId string", err)
			return nil, apierrors.NewInternalServerApiError("error when tying to get all items", error_utils.GetDatabaseGenericError())
		}
		result = append(result, itemId)
	}

	return result, nil
}
