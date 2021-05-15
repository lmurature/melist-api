package share

import (
	"fmt"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
)

const (
	ShareTypeRead = "read"
	ShareTypeWrite = "write"
	ShareTypeCheck = "check"
)

type ShareConfig struct {
	ListId    int64  `json:"list_id"`
	UserId    int64  `json:"user_id"`
	ShareType string `json:"share_type"`
}

type ShareConfigs []ShareConfig

func (s ShareConfig) Validate() apierrors.ApiError {
	if s.ShareType == "" {
		return apierrors.NewBadRequestApiError(fmt.Sprintf("share type cant be empty for user id %d", s.UserId))
	}

	if s.ShareType != ShareTypeWrite && s.ShareType != ShareTypeRead && s.ShareType  != ShareTypeCheck {
		return apierrors.NewBadRequestApiError(fmt.Sprintf("share type must be 'read' 'write' or 'check' for user id %d", s.UserId))
	}

	return nil
}
