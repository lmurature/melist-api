package lists

import (
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	"github.com/lmurature/melist-api/src/api/domain/lists"
	"github.com/lmurature/melist-api/src/api/domain/share"
)

type listsService struct{}

func (l listsService) CreateList(dto lists.ListDto) (lists.ListDto, apierrors.ApiError) {
	panic("implement me")
}

func (l listsService) UpdateList(dto lists.ListDto) (lists.ListDto, apierrors.ApiError) {
	panic("implement me")
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

type listsServiceInterface interface {
	CreateList(dto lists.ListDto) (lists.ListDto, apierrors.ApiError)
	UpdateList(dto lists.ListDto) (lists.ListDto, apierrors.ApiError)
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