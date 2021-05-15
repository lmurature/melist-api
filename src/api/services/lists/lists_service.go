package lists

import (
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	"github.com/lmurature/melist-api/src/api/domain/lists"
	"github.com/lmurature/melist-api/src/api/domain/share"
	date_utils "github.com/lmurature/melist-api/src/api/utils/date"
	"github.com/sirupsen/logrus"
	"net/http"
)

type listsService struct{}

type listsServiceInterface interface {
	CreateList(dto lists.ListDto) (*lists.ListDto, apierrors.ApiError)
	UpdateList(dto lists.ListDto, callerId int64) (*lists.ListDto, apierrors.ApiError)
	GetList(listId int64, callerId int64) (*lists.ListDto, apierrors.ApiError)
	GiveAccessToUsers(listId int64, callerId int64, config []share.ShareConfig) apierrors.ApiError
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

func (l listsService) UpdateList(updatedList lists.ListDto, callerId int64) (*lists.ListDto, apierrors.ApiError) {
	actualList := lists.ListDto{Id: updatedList.Id}
	if err := actualList.Get(); err != nil {
		return nil, err
	}

	if err := actualList.ValidateUpdatability(callerId); err != nil {
		return nil, err
	}

	actualList.UpdateFields(updatedList)

	if err := actualList.Update(); err != nil {
		return nil, err
	}

	return &actualList, nil
}

func (l listsService) GetList(listId int64, callerId int64) (*lists.ListDto, apierrors.ApiError) {
	dto := lists.ListDto{Id: listId}

	if err := dto.Get(); err != nil {
		return nil, err
	}

	listShareConfigs := share.ShareConfig{ListId: listId}
	configs, err := listShareConfigs.GetAllSharedConfigsByList()
	if err != nil {
		return nil, err
	}

	if err := dto.ValidateReadability(callerId, configs); err != nil {
		return nil, err
	}

	return &dto, nil
}

func (l listsService) GiveAccessToUsers(listId int64, callerId int64, config []share.ShareConfig) apierrors.ApiError {
	list := lists.ListDto{Id: listId}
	if err := list.Get(); err != nil {
		return err
	}

	if err := list.ValidateUpdatability(callerId); err != nil {
		return err
	}

	errorCauseList := make(apierrors.CauseList, 0)
	for i, c := range config {
		if err := c.Validate(); err != nil {
			errorCauseList = append(errorCauseList, err.Message())
		}
		config[i].ListId = listId
	}

	if len(errorCauseList) > 0 {
		return apierrors.NewApiError("invalid request", "invalid request for sharing access to users", http.StatusBadRequest, errorCauseList)
	}

	errorSaving := make(apierrors.CauseList, 0)
	for i := range config {
		if err := config[i].Save(); err != nil {
			logrus.Error("error while trying to save share config", err)
			errorSaving = append(errorSaving, err.Message())
		}
	}

	if len(errorSaving) > 0 {
		return apierrors.NewApiError("error while trying to save share configs", "database error", http.StatusInternalServerError, errorSaving)
	}

	return nil
}

func (l listsService) SearchPublicLists() (lists.Lists, apierrors.ApiError) {
	dto := lists.ListDto{}
	return dto.GetAllPublicLists()
}

func (l listsService) GetMyLists(ownerId int64) (lists.Lists, apierrors.ApiError) {
	dto := lists.ListDto{OwnerId: ownerId}
	return dto.GetListsFromOwner()
}

func (l listsService) GetMySharedLists(userId int64) (lists.Lists, apierrors.ApiError) {
	dto := share.ShareConfig{UserId: userId}
	userSharedConfigs, err := dto.GetAllShareConfigsByUser()
	if err != nil {
		return nil, err
	}

	result := make([]lists.ListDto, 0)
	for _, c := range userSharedConfigs {
		listDto := lists.ListDto{Id: c.ListId}
		if err := listDto.Get(); err != nil {
			return nil, err
		}
		result = append(result, listDto)
	}

	if len(result) == 0 {
		return nil, apierrors.NewNotFoundApiError("you have no shared lists")
	}

	return result, nil
}
