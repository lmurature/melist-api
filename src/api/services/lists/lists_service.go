package lists

import (
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	"github.com/lmurature/melist-api/src/api/domain/lists"
	"github.com/lmurature/melist-api/src/api/domain/share"
	date_utils "github.com/lmurature/melist-api/src/api/utils/date"
)

type listsService struct{}

type listsServiceInterface interface {
	CreateList(dto lists.ListDto) (*lists.ListDto, apierrors.ApiError)
	UpdateList(dto lists.ListDto) (*lists.ListDto, apierrors.ApiError)
	GetList(listId int64, callerId int64) (*lists.ListDto, apierrors.ApiError)
	GiveAccessToUsers(config []share.ShareConfig) apierrors.ApiError
	SearchPublicLists() (lists.Lists, apierrors.ApiError)
	GetMyLists(ownerId int64) (lists.Lists, apierrors.ApiError)
	GetMySharedLists(userId int64) (lists.Lists, apierrors.ApiError)
}

var (
	ListsService listsServiceInterface
)

func init() {
	ListsService = listsService{}
}

func (l listsService) CreateList(dto lists.ListDto) (*lists.ListDto, apierrors.ApiError) {
	if err := dto.Validate(); err != nil {
		return nil, err
	}

	dto.DateCreated = date_utils.GetNowDateFormatted()
	if err := dto.Save(); err != nil {
		return nil, err
	}

	return &dto, nil
}

func (l listsService) UpdateList(dto lists.ListDto) (*lists.ListDto, apierrors.ApiError) {
	panic("implement me")
}

func (l listsService) GetList(listId int64, callerId int64) (*lists.ListDto, apierrors.ApiError) {
	dto := lists.ListDto{Id: listId}

	if err := dto.Get(); err != nil {
		return nil, err
	}

	if err := dto.ValidateReadability(callerId); err != nil {
		return nil, err
	}

	return &dto, nil
}

func (l listsService) GiveAccessToUsers(config []share.ShareConfig) apierrors.ApiError {
	panic("implement me")
}

func (l listsService) SearchPublicLists() (lists.Lists, apierrors.ApiError) {
	panic("implement me")
}

func (l listsService) GetMyLists(ownerId int64) (lists.Lists, apierrors.ApiError) {
	panic("implement me")
}

func (l listsService) GetMySharedLists(userId int64) (lists.Lists, apierrors.ApiError) {
	panic("implement me")
}
