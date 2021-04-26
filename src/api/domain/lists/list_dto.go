package lists

import "github.com/lmurature/melist-api/src/api/domain/apierrors"

const (
	PrivacyTypePrivate = "private"
	PrivacyTypePublic  = "public"
)

type ListDto struct {
	Id          int64  `json:"id"`
	OwnerId     int64  `json:"owner_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
	DateCreated string `json:"date_created"`
}

type Lists []ListDto

func (l ListDto) Validate() apierrors.ApiError {
	if l.OwnerId == 0 || l.Title == "" || l.Privacy == "" || (l.Privacy != PrivacyTypePrivate && l.Privacy != PrivacyTypePublic) {
		return apierrors.NewBadRequestApiError("invalid list values. Required values: title, privacy (private or public)")
	}
	return nil
}

func (l ListDto) ValidateUpdatability(callerId int64) apierrors.ApiError {
	if l.OwnerId != callerId {
		return apierrors.NewForbiddenApiError("you cannot update this list")
	}
	return nil
}

func (l ListDto) ValidateReadability(callerId int64) apierrors.ApiError {
	if l.Privacy == PrivacyTypePrivate {
		if callerId != l.OwnerId {
			return apierrors.NewForbiddenApiError("you have no access to this list")
		}
	}
	return nil
}

func (l *ListDto) UpdateFields(updatedList ListDto) {
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

