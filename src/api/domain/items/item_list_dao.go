package items

import (
	"fmt"
	"github.com/lmurature/melist-api/src/api/clients/database"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	error_utils "github.com/lmurature/melist-api/src/api/utils/error"
	"github.com/sirupsen/logrus"
)

const (
	getItemsFromList   = "SELECT l.list_id, l.item_id, l.status, l.variation_external_id FROM list_item l WHERE l.list_id=?;"
	insertItemToList   = "INSERT INTO list_item(list_id, item_id, status, variation_external_id) VALUES(?,?,?,?);"
	removeItemFromList = "DELETE FROM list_item l WHERE l.item_id=? and l.list_id=?;"
	checkItem          = "UPDATE list_item SET status=? WHERE item_id=? and list_id=?;"
)

var (
	ItemListDao itemListDaoInterface
)

type itemListDaoInterface interface {
	InsertItemToList(itemList ItemListDto) (*ItemListDto, apierrors.ApiError)
	DeleteItemFromList(itemId string, listId int64) apierrors.ApiError
	GetItemsFromList(listId int64) (ItemListCollection, apierrors.ApiError)
	UpdateItemStatus(itemId string, listId int64, status string) apierrors.ApiError
}

type itemListDao struct{}

func init() {
	ItemListDao = &itemListDao{}
}

func (dao *itemListDao) InsertItemToList(itemList ItemListDto) (*ItemListDto, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(insertItemToList)
	if err != nil {
		logrus.Error("error when trying to prepare insert item to list statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to insert item to list", error_utils.GetDatabaseGenericError())
	}
	defer stmt.Close()

	_, saveErr := stmt.Exec(itemList.ListId, itemList.ItemId, itemList.Status, itemList.VariationId)
	if saveErr != nil {
		return nil, apierrors.NewInternalServerApiError("error when trying to save item to list", error_utils.GetDatabaseGenericError())
	}

	logrus.Info(fmt.Sprintf("successfully added item %s to list %d", itemList.ItemId, itemList.ListId))
	return &itemList, nil
}

func (dao *itemListDao) DeleteItemFromList(itemId string, listId int64) apierrors.ApiError {
	stmt, err := database.DbClient.Prepare(removeItemFromList)
	if err != nil {
		logrus.Error("error when trying to prepare delete item from list statement", err)
		return apierrors.NewInternalServerApiError("error when trying to delete item from list", error_utils.GetDatabaseGenericError())
	}
	defer stmt.Close()

	_, deleteErr := stmt.Exec(itemId, listId)
	if deleteErr != nil {
		return apierrors.NewInternalServerApiError("error when trying to delete item from list", error_utils.GetDatabaseGenericError())
	}

	logrus.Info(fmt.Sprintf("successfully deleted item %s from list %d", itemId, listId))
	return nil
}

func (dao *itemListDao) GetItemsFromList(listId int64) (ItemListCollection, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(getItemsFromList)
	if err != nil {
		logrus.Error("error when trying to prepare get all items from list statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to get all items from list", error_utils.GetDatabaseGenericError())
	}
	defer stmt.Close()

	rows, err := stmt.Query(listId)
	if err != nil {
		logrus.Error("error while getting all items from list", err)
		return nil, apierrors.NewInternalServerApiError("error getting all items from list", error_utils.GetDatabaseGenericError())
	}
	defer rows.Close()

	result := make([]ItemListDto, 0)

	for rows.Next() {
		var i ItemListDto
		if err := rows.Scan(&i.ListId, &i.ItemId, &i.Status, &i.VariationId); err != nil {
			logrus.Error("error when scan item row into item struct", err)
			return nil, apierrors.NewInternalServerApiError("error when tying to get all items from list", error_utils.GetDatabaseGenericError())
		}
		result = append(result, i)
	}

	return result, nil
}

func (dao *itemListDao) UpdateItemStatus(itemId string, listId int64, status string) apierrors.ApiError {
	stmt, err := database.DbClient.Prepare(checkItem)
	if err != nil {
		logrus.Error("error when trying to prepare get check item from list statement", err)
		return apierrors.NewInternalServerApiError("error when trying to item from list", error_utils.GetDatabaseGenericError())
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(status, itemId, listId)
	if updateErr != nil {
		logrus.Error("error when trying to check item", updateErr)
		return apierrors.NewInternalServerApiError("error when trying to check item", error_utils.GetDatabaseGenericError())
	}

	logrus.Info(fmt.Sprintf("successfully updated item %s from list %d", itemId, listId))
	return nil
}
