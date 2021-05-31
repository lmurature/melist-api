package items

import (
	"fmt"
	"github.com/lmurature/melist-api/src/api/clients/database"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	error_utils "github.com/lmurature/melist-api/src/api/utils/error"
	"github.com/sirupsen/logrus"
)

const(
	insertItem = "INSERT INTO item(item_id) VALUES(?);"
)

type ItemDto string

func (i *ItemDto) InsertItem() apierrors.ApiError {
	stmt, err := database.DbClient.Prepare(insertItem)
	if err != nil {
		logrus.Error("error when trying to prepare insert item statement", err)
		return apierrors.NewInternalServerApiError("error when trying to insert item", error_utils.GetDatabaseGenericError())
	}
	defer stmt.Close()

	_, saveErr := stmt.Exec(i)
	if saveErr != nil {
		return apierrors.NewInternalServerApiError("error when trying to save item to items table", error_utils.GetDatabaseGenericError())
	}

	logrus.Info(fmt.Sprintf("successfully added %s to item table", *i))
	return nil
}
