package error_utils

import "errors"

func GetDatabaseGenericError() error {
	return errors.New("database error")
}
