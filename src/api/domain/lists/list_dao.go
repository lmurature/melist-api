package lists

import (
	"errors"
	"fmt"
	"github.com/lmurature/melist-api/src/api/clients/database"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	"github.com/sirupsen/logrus"
)

const (
	getList = "SELECT l.id, l.owner_id, l.title, l.description, l.privacy, l.date_created FROM list l WHERE l.id=?;"
	insertList = "INSERT INTO list(owner_id, title, description, privacy, date_created) VALUES(?,?,?,?,?);"
	updateList = "UPDATE list SET title=?, description=?, privacy=? WHERE id=?;"
	getAllPublicLists = "SELECT l.id, l.owner_id, l.title, l.description, l.privacy, l.date_created FROM list l WHERE l.privacy='public';"
	getAllListsFromOwner = "SELECT l.id, l.owner_id, l.title, l.description, l.privacy, l.date_created FROM list l WHERE l.id=?;"
)

func (l *ListDto) Get() apierrors.ApiError {
	stmt, err := database.DbClient.Prepare(getList)
	if err != nil {
		logrus.Error("error when trying to prepare get list statement", err)
		return apierrors.NewInternalServerApiError("error when trying to get list", errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(l.Id)

	if queryErr := result.Scan(l.Id, l.OwnerId, l.Title, l.Description, l.Privacy, l.DateCreated); queryErr != nil {
		msg := fmt.Sprintf("list %d not found", l.Id)
		logrus.Error(msg)
		return apierrors.NewNotFoundApiError(msg)
	}

	return nil
}

func (l *ListDto) Save() apierrors.ApiError {
	stmt, err := database.DbClient.Prepare(insertList)
	if err != nil {
		logrus.Error("error when trying to prepare insert list statement", err)
		return apierrors.NewInternalServerApiError("error when trying to insert list", errors.New("database error"))
	}
	defer stmt.Close()

	execResult, execErr := stmt.Exec(l.OwnerId, l.Title, l.Description, l.Privacy, l.DateCreated)
	if execErr != nil {
		logrus.Error("error when trying to save list", execErr)
		return apierrors.NewInternalServerApiError("error when trying to save list", errors.New("database error"))
	}

	listId, _ := execResult.LastInsertId()

	logrus.Info(fmt.Sprintf("successfully created list %d", listId))
	l.Id = listId

	return nil
}

func (l *ListDto) Update() apierrors.ApiError {
	stmt, err := database.DbClient.Prepare(updateList)
	if err != nil {
		logrus.Error("error when trying to prepare update list statement", err)
		return apierrors.NewInternalServerApiError("error when update to insert list", errors.New("database error"))
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(l.Title, l.Description, l.Privacy, l.Id)
	if updateErr != nil {
		logrus.Error("error when trying to update list", updateErr)
		return apierrors.NewInternalServerApiError("error when trying to update list", errors.New("database error"))
	}

	logrus.Info(fmt.Sprintf("successfully updated list %d", l.Id))
	return nil
}

func (l *ListDto) GetAllPublicLists() (Lists, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(getAllPublicLists)
	if err != nil {
		logrus.Error("error when trying to prepare get public lists statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to get public lists", errors.New("database error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		logrus.Error("error while getting public lists", err)
		return nil, apierrors.NewInternalServerApiError("error getting public lists", errors.New("database error"))
	}
	defer rows.Close()

	result := make([]ListDto, 0)

	for rows.Next() {
		var list ListDto
		if err := rows.Scan(&list.Id, &list.OwnerId, &list.Title, &list.Description, &list.Privacy, &list.DateCreated); err != nil {
			logrus.Error("error when scan list row into list struct", err)
			return nil, apierrors.NewInternalServerApiError("error when tying to get public lists", errors.New("database error"))
		}
		result = append(result, list)
	}

	if len(result) == 0 {
		return nil, apierrors.NewNotFoundApiError("no public lists found")
	}

	return result, nil
}

func (l *ListDto) GetListsFromOwner() (Lists, apierrors.ApiError) {
	stmt, err := database.DbClient.Prepare(getAllListsFromOwner)
	if err != nil {
		logrus.Error("error when trying to prepare get all owner lists statement", err)
		return nil, apierrors.NewInternalServerApiError("error when trying to get all owner lists", errors.New("database error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(l.OwnerId)
	if err != nil {
		logrus.Error("error while getting all owner lists", err)
		return nil, apierrors.NewInternalServerApiError("error getting all owner lists", errors.New("database error"))
	}
	defer rows.Close()

	result := make([]ListDto, 0)

	for rows.Next() {
		var list ListDto
		if err := rows.Scan(&list.Id, &list.OwnerId, &list.Title, &list.Description, &list.Privacy, &list.DateCreated); err != nil {
			logrus.Error("error when scan list row into list struct", err)
			return nil, apierrors.NewInternalServerApiError("error when tying to get all owner lists", errors.New("database error"))
		}
		result = append(result, list)
	}

	if len(result) == 0 {
		return nil, apierrors.NewNotFoundApiError(fmt.Sprintf("no lists found for owner %d", l.OwnerId))
	}

	return result, nil
}
