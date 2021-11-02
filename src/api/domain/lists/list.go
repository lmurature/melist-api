package lists

import (
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	"github.com/lmurature/melist-api/src/api/domain/share"
)

const (
	PrivacyTypePrivate = "private"
	PrivacyTypePublic  = "public"
)

type List struct {
	Id            int64  `json:"id"`
	OwnerId       int64  `json:"owner_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Privacy       string `json:"privacy"`
	DateCreated   string `json:"date_created"`
	Notifications int    `json:"notifications,omitempty"`
}

type Lists []List

func (lists Lists) ContainsList(id int64) bool {
	for _, l := range lists {
		if l.Id == id {
			return true
		}
	}
	return false
}

func (l List) Validate() apierrors.ApiError {
	if l.OwnerId == 0 || l.Title == "" || l.Privacy == "" || (l.Privacy != PrivacyTypePrivate && l.Privacy != PrivacyTypePublic) {
		return apierrors.NewBadRequestApiError("invalid list values. Required values: title, privacy (private or public)")
	}
	return nil
}

func (l List) ValidateAddItems(callerId int64, configs share.ShareConfigs) apierrors.ApiError {
	if l.Privacy == PrivacyTypePrivate {
		if callerId == l.OwnerId {
			return nil
		}
		for _, c := range configs {
			if c.UserId == callerId && c.ShareType == share.ShareTypeWrite {
				return nil
			}
		}
	} else {
		return nil
	}
	return apierrors.NewForbiddenApiError("you have no access to this list")
}

func (l List) ValidateUpdatability(callerId int64) apierrors.ApiError {
	if l.OwnerId != callerId {
		return apierrors.NewForbiddenApiError("you have no access to update this list")
	}
	return nil
}

func (l List) ValidateCheckItems(callerId int64, configs share.ShareConfigs) apierrors.ApiError {
	if l.Privacy == PrivacyTypePrivate {
		if callerId == l.OwnerId {
			return nil
		}
		for _, c := range configs {
			if c.UserId == callerId && (c.ShareType == share.ShareTypeCheck || c.ShareType == share.ShareTypeWrite) {
				return nil
			}
		}
	} else {
		return nil
	}
	return apierrors.NewForbiddenApiError("you have no access to this list")
}

func (l List) ValidateReadability(callerId int64, configs share.ShareConfigs) apierrors.ApiError {
	if l.Privacy == PrivacyTypePrivate {
		if callerId == l.OwnerId {
			return nil
		}
		for _, c := range configs {
			if c.UserId == callerId {
				return nil
			}
		}
	} else {
		return nil
	}
	return apierrors.NewForbiddenApiError("you have no access to this list")
}

func (l *List) UpdateFields(updatedList List) {
	if updatedList.Privacy != "" && (updatedList.Privacy == PrivacyTypePrivate || updatedList.Privacy == PrivacyTypePublic) {
		l.Privacy = updatedList.Privacy
	}
	if updatedList.Title != "" {
		l.Title = updatedList.Title
	}
	if updatedList.Description != "" {
		l.Description = updatedList.Description
	}
}
