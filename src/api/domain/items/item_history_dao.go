package items

import (
	"github.com/lmurature/melist-api/src/api/clients/database"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	error_utils "github.com/lmurature/melist-api/src/api/utils/error"
	"github.com/sirupsen/logrus"
)

const (
	insertItemHistory  = "INSERT INTO item_history(item_id,price,quantity,status,has_deal,date_fetched) VALUES(?,?,?,?,?,?);"
	getItemHistory     = "SELECT id,item_id,price,quantity,status,has_deal,date_fetched FROM item_history WHERE item_id=? ORDER BY date_fetched ASC;"
	getLastItemHistory = "SELECT id,item_id,price,quantity,status,has_deal,date_fetched FROM item_history WHERE item_id=? ORDER BY date_fetched DESC LIMIT 1;"
)

var (
	ItemHistoryDao itemHistoryDaoInterface
)

type itemHistoryDaoInterface interface {
	InsertItemHistory(itemHistory ItemHistory) (*ItemHistory, apierrors.ApiError)
	GetLastItemHistory(itemId string) (*ItemHistory, apierrors.ApiError)
	GetItemHistory(itemId string) ([]ItemHistory, apierrors.ApiError)
}

type itemHistoryDao struct{}

func init() {
	ItemHistoryDao = &itemHistoryDao{}
}

func (dao *itemHistoryDao) InsertItemHistory(history ItemHistory) (*ItemHistory, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(insertItemHistory)
	if err != nil {
		logrus.Error("error when trying to prepare insert item history statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to insert item history", error_utils.GetDatabaseGenericError())
	}
	defer stmt.Close()

	execResult, execErr := stmt.Exec(history.ItemId, history.Price, history.Quantity, history.Status, history.HasDeal, history.DateFetched)
	if execErr != nil {
		logrus.Error("error when trying to save item history", execErr)
		return nil, apierrors.NewInternalServerApiError("error when trying to save item history", error_utils.GetDatabaseGenericError())
	}

	historyId, _ := execResult.LastInsertId()

	history.Id = historyId

	return &history, nil
}

func (dao *itemHistoryDao) GetLastItemHistory(itemId string) (*ItemHistory, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(getLastItemHistory)
	if err != nil {
		logrus.Error("error when trying to prepare get item history statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to get item history", error_utils.GetDatabaseGenericError())
	}
	defer stmt.Close()

	rows, err := stmt.Query(itemId)
	if err != nil {
		logrus.Error("error while getting item history", err)
		return nil, apierrors.NewInternalServerApiError("error getting item history", error_utils.GetDatabaseGenericError())
	}
	defer rows.Close()

	var result *ItemHistory = nil
	for rows.Next() {
		var history ItemHistory
		if err := rows.Scan(&history.Id, &history.ItemId, &history.Price,
			&history.Quantity, &history.Status, &history.HasDeal,
			&history.DateFetched); err != nil {
			logrus.Error("error scaning row item history", err)
			return nil, apierrors.NewInternalServerApiError("error scanning row item history", error_utils.GetDatabaseGenericError())
		}
		result = &history
	}

	if result == nil {
		return nil, apierrors.NewNotFoundApiError("item history not found")
	}

	return result, nil
}

func (dao *itemHistoryDao) GetItemHistory(itemId string) ([]ItemHistory, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(getItemHistory)
	if err != nil {
		logrus.Error("error when trying to prepare get item history statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to get item history", error_utils.GetDatabaseGenericError())
	}
	defer stmt.Close()

	rows, err := stmt.Query(itemId)
	if err != nil {
		logrus.Error("error while getting item history", err)
		return nil, apierrors.NewInternalServerApiError("error getting item history", error_utils.GetDatabaseGenericError())
	}
	defer rows.Close()

	result := make([]ItemHistory, 0)
	for rows.Next() {
		var history ItemHistory
		if err := rows.Scan(&history.Id, &history.ItemId, &history.Price,
			&history.Quantity, &history.Status, &history.HasDeal,
			&history.DateFetched); err != nil {
			logrus.Error("error scaning row item history", err)
			return nil, apierrors.NewInternalServerApiError("error scanning row item history", error_utils.GetDatabaseGenericError())
		}
		result = append(result, history)
	}

	return result, nil
}
