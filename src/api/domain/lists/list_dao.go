package lists

import (
	"fmt"
	"github.com/lmurature/melist-api/src/api/clients/database"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	error_utils "github.com/lmurature/melist-api/src/api/utils/error"
	"github.com/sirupsen/logrus"
)

const (
	getList              = "SELECT l.id, l.owner_id, l.title, l.description, l.privacy, l.date_created FROM list l WHERE l.id=?;"
	insertList           = "INSERT INTO list(owner_id, title, description, privacy, date_created) VALUES(?,?,?,?,?);"
	updateList           = "UPDATE list SET title=?, description=?, privacy=? WHERE id=?;"
	getAllPublicLists    = "SELECT l.id, l.owner_id, l.title, l.description, l.privacy, l.date_created FROM list l WHERE l.privacy='public';"
	getAllListsFromOwner = "SELECT l.id, l.owner_id, l.title, l.description, l.privacy, l.date_created FROM list l WHERE l.owner_id=?;"
)

var (
	ListDao listDaoInterface
)

type listDaoInterface interface {
	GetList(listId int64) (*List, apierrors.ApiError)
	CreateList(listDto List) (*List, apierrors.ApiError)
	UpdateList(listDto List) (*List, apierrors.ApiError)
	GetPublicLists() (Lists, apierrors.ApiError)
	GetListsFromOwner(ownerId int64) (Lists, apierrors.ApiError)
}

type listDao struct{}

func init() {
	ListDao = &listDao{}
}

func (dao *listDao) GetList(listId int64) (*List, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(getList)
	if err != nil {
		logrus.Error("error when trying to prepare get list statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to get list", error_utils.GetDatabaseGenericError())
	}
	defer stmt.Close()

	result := stmt.QueryRow(listId)

	var listDto List
	if queryErr := result.Scan(&listDto.Id, &listDto.OwnerId, &listDto.Title, &listDto.Description,
		&listDto.Privacy, &listDto.DateCreated); queryErr != nil {
		msg := fmt.Sprintf("list %d not found", listId)
		logrus.Error(msg, queryErr)
		return nil, apierrors.NewNotFoundApiError(msg)
	}

	return &listDto, nil
}

func (dao *listDao) CreateList(listDto List) (*List, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(insertList)
	if err != nil {
		logrus.Error("error when trying to prepare insert list statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to insert list", error_utils.GetDatabaseGenericError())
	}
	defer stmt.Close()

	execResult, execErr := stmt.Exec(listDto.OwnerId, listDto.Title, listDto.Description, listDto.Privacy, listDto.DateCreated)
	if execErr != nil {
		logrus.Error("error when trying to save list", execErr)
		return nil, apierrors.NewInternalServerApiError("error when trying to save list", error_utils.GetDatabaseGenericError())
	}

	listId, _ := execResult.LastInsertId()

	logrus.Info(fmt.Sprintf("successfully created list %d", listId))
	listDto.Id = listId

	return &listDto, nil
}

func (dao *listDao) UpdateList(listDto List) (*List, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(updateList)
	if err != nil {
		logrus.Error("error when trying to prepare update list statement", err)
		return nil, apierrors.NewInternalServerApiError("error when update to insert list", error_utils.GetDatabaseGenericError())
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(listDto.Title, listDto.Description, listDto.Privacy, listDto.Id)
	if updateErr != nil {
		logrus.Error("error when trying to update list", updateErr)
		return nil, apierrors.NewInternalServerApiError("error when trying to update list", error_utils.GetDatabaseGenericError())
	}

	logrus.Info(fmt.Sprintf("successfully updated list %d", listDto.Id))
	return &listDto, nil
}

func (dao *listDao) GetPublicLists() (Lists, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(getAllPublicLists)
	if err != nil {
		logrus.Error("error when trying to prepare get public lists statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to get public lists", error_utils.GetDatabaseGenericError())
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		logrus.Error("error while getting public lists", err)
		return nil, apierrors.NewInternalServerApiError("error getting public lists", error_utils.GetDatabaseGenericError())
	}
	defer rows.Close()

	result := make([]List, 0)

	for rows.Next() {
		var list List
		if err := rows.Scan(&list.Id, &list.OwnerId, &list.Title, &list.Description, &list.Privacy, &list.DateCreated); err != nil {
			logrus.Error("error when scan list row into list struct", err)
			return nil, apierrors.NewInternalServerApiError("error when tying to get public lists", error_utils.GetDatabaseGenericError())
		}
		result = append(result, list)
	}

	if len(result) == 0 {
		return nil, apierrors.NewNotFoundApiError("no public lists found")
	}

	return result, nil
}

func (dao *listDao) GetListsFromOwner(ownerId int64) (Lists, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(getAllListsFromOwner)
	if err != nil {
		logrus.Error("error when trying to prepare get all owner lists statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to get all owner lists", error_utils.GetDatabaseGenericError())
	}
	defer stmt.Close()

	rows, err := stmt.Query(ownerId)
	if err != nil {
		logrus.Error("error while getting all owner lists", err)
		return nil, apierrors.NewInternalServerApiError("error getting all owner lists", error_utils.GetDatabaseGenericError())
	}
	defer rows.Close()

	result := make([]List, 0)

	for rows.Next() {
		var list List
		if err := rows.Scan(&list.Id, &list.OwnerId, &list.Title, &list.Description, &list.Privacy, &list.DateCreated); err != nil {
			logrus.Error("error when scan list row into list struct", err)
			return nil, apierrors.NewInternalServerApiError("error when tying to get all owner lists", error_utils.GetDatabaseGenericError())
		}
		result = append(result, list)
	}

	if len(result) == 0 {
		return nil, apierrors.NewNotFoundApiError(fmt.Sprintf("no lists found for owner %d", ownerId))
	}

	return result, nil
}
