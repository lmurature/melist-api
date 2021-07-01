package lists

import (
	"fmt"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	"github.com/lmurature/melist-api/src/api/domain/items"
	"github.com/lmurature/melist-api/src/api/domain/lists"
	"github.com/lmurature/melist-api/src/api/domain/share"
	items_service "github.com/lmurature/melist-api/src/api/services/items"
	date_utils "github.com/lmurature/melist-api/src/api/utils/date"
	"github.com/lmurature/melist-api/src/api/utils/slice"
	"github.com/sirupsen/logrus"
	"net/http"
)

type listsService struct{}

type listsServiceInterface interface {
	CreateList(dto lists.List) (*lists.List, apierrors.ApiError)
	UpdateList(dto lists.List, callerId int64) (*lists.List, apierrors.ApiError)
	GetList(listId int64, callerId int64) (*lists.List, apierrors.ApiError)
	GetListShareConfigs(listId int64, callerId int64) (share.ShareConfigs, apierrors.ApiError)
	GiveAccessToUsers(listId int64, callerId int64, config share.ShareConfigs) apierrors.ApiError
	SearchPublicLists() (lists.Lists, apierrors.ApiError)
	GetMyLists(ownerId int64) (lists.Lists, apierrors.ApiError)
	GetMySharedLists(userId int64) (lists.Lists, apierrors.ApiError)
	AddItemToList(itemId string, variationId int64, listId int64, callerId int64) apierrors.ApiError
	GetItemsFromList(listId int64, callerId int64, info bool) (items.ItemListCollection, apierrors.ApiError)
	DeleteItemFromList(itemId string, listId int64, callerId int64) apierrors.ApiError
	CheckItem(itemId string, listId int64, callerId int64) apierrors.ApiError
}

var (
	ListsService listsServiceInterface
)

func init() {
	ListsService = listsService{}
}

func (l listsService) CreateList(list lists.List) (*lists.List, apierrors.ApiError) {
	if err := list.Validate(); err != nil {
		return nil, err
	}

	list.DateCreated = date_utils.GetNowDateFormatted()

	result, err := lists.ListDao.CreateList(list)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (l listsService) UpdateList(updatedList lists.List, callerId int64) (*lists.List, apierrors.ApiError) {
	actualList, err := lists.ListDao.GetList(updatedList.Id)
	if err != nil {
		return nil, err
	}

	if err := actualList.ValidateUpdatability(callerId); err != nil {
		return nil, err
	}

	actualList.UpdateFields(updatedList)

	result, err := lists.ListDao.UpdateList(*actualList)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (l listsService) GetList(listId int64, callerId int64) (*lists.List, apierrors.ApiError) {
	list, err := lists.ListDao.GetList(listId)
	if err != nil {
		return nil, err
	}

	configs, err := share.ShareConfigDao.GetAllShareConfigsByList(listId)
	if err != nil {
		if err.Status() != http.StatusNotFound {
			return nil, err
		}
	}

	if err := list.ValidateReadability(callerId, configs); err != nil {
		return nil, err
	}

	return list, nil
}

func (l listsService) GiveAccessToUsers(listId int64, callerId int64, config share.ShareConfigs) apierrors.ApiError {
	list, err := lists.ListDao.GetList(listId)
	if err != nil {
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

	actualConfigs, err := share.ShareConfigDao.GetAllShareConfigsByList(listId)
	if err != nil {
		if err.Status() != http.StatusNotFound {
			return err
		}
	}

	errorSaving := make(apierrors.CauseList, 0)
	for i := range config {

		var dbErr apierrors.ApiError
		if slice.ShareConfigUserExists(actualConfigs, config[i].UserId) {
			_, dbErr = share.ShareConfigDao.UpdateShareConfig(config[i])
		} else {
			_, dbErr = share.ShareConfigDao.CreateShareConfig(config[i])
		}

		if dbErr != nil {
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
	return lists.ListDao.GetPublicLists()
}

func (l listsService) GetMyLists(ownerId int64) (lists.Lists, apierrors.ApiError) {
	return lists.ListDao.GetListsFromOwner(ownerId)
}

func (l listsService) GetMySharedLists(userId int64) (lists.Lists, apierrors.ApiError) {
	userSharedConfigs, err := share.ShareConfigDao.GetAllShareConfigsByUser(userId)
	if err != nil {
		return nil, err
	}

	result := make([]lists.List, 0)
	for _, c := range userSharedConfigs {
		listDto, err := lists.ListDao.GetList(c.ListId)
		if err != nil {
			return nil, err
		}
		result = append(result, *listDto)
	}

	if len(result) == 0 {
		return nil, apierrors.NewNotFoundApiError("you have no shared lists")
	}

	return result, nil
}

func (l listsService) GetListShareConfigs(listId int64, callerId int64) (share.ShareConfigs, apierrors.ApiError) {
	list, err := lists.ListDao.GetList(listId)
	if err != nil {
		return nil, err
	}

	if list.OwnerId != callerId {
		return nil, apierrors.NewUnauthorizedApiError("you have no access to this list's share config")
	}

	configList, err := share.ShareConfigDao.GetAllShareConfigsByList(listId)
	if err != nil {
		return nil, err
	}

	return configList, nil
}

func (l listsService) AddItemToList(itemId string, variationId int64, listId int64, callerId int64) apierrors.ApiError {
	list, err := lists.ListDao.GetList(listId)
	if err != nil {
		return err
	}

	actualConfigs, err := share.ShareConfigDao.GetAllShareConfigsByList(listId)
	if err != nil {
		if err.Status() != http.StatusNotFound {
			return err
		}
	}

	if err := list.ValidateAddItems(callerId, actualConfigs); err != nil {
		return err
	}

	// Check if list already has the item
	itemCollection, err := items.ItemListDao.GetItemsFromList(listId)
	if err != nil {
		return err
	}

	if itemCollection.ContainsItem(itemId) {
		return apierrors.NewBadRequestApiError(fmt.Sprintf("item %s is already in the list", itemId))
	}

	// insert into item table
	if err := items.ItemDao.InsertItem(itemId); err != nil {
		return err
	}

	itemListDto := items.ItemListDto{
		ItemId:      itemId,
		ListId:      listId,
		Status:      items.StatusNotChecked,
		VariationId: variationId,
	}

	_, err = items.ItemListDao.InsertItemToList(itemListDto)
	if err != nil {
		return err
	}

	return nil
}

func (l listsService) GetItemsFromList(listId int64, callerId int64, info bool) (items.ItemListCollection, apierrors.ApiError) {
	list, err := lists.ListDao.GetList(listId)
	if err != nil {
		return nil, err
	}

	actualConfigs, err := share.ShareConfigDao.GetAllShareConfigsByList(listId)
	if err != nil {
		if err.Status() != http.StatusNotFound {
			return nil, err
		}
	}

	if err := list.ValidateReadability(callerId, actualConfigs); err != nil {
		return nil, err
	}

	itemListCollection, err := items.ItemListDao.GetItemsFromList(listId)
	if err != nil {
		return nil, err
	}

	if info {
		input := make(chan items.ItemConcurrent, len(itemListCollection))
		defer close(input)

		for i := range itemListCollection {
			go func(id string, index int, output chan items.ItemConcurrent) {
				item, err := items_service.ItemsService.GetItem(id)
				output <- items.ItemConcurrent{
					Item:      item,
					Error:     err,
					ListIndex: index,
				}
			}(itemListCollection[i].ItemId, i, input)
		}

		for i := 0; i < len(itemListCollection); i++ {
			result := <-input
			if result.Error != nil {
				return nil, result.Error
			}

			itemListCollection[result.ListIndex].MeliItem = result.Item
		}
	}

	return itemListCollection, nil
}

func (l listsService) DeleteItemFromList(itemId string, listId int64, callerId int64) apierrors.ApiError {
	list, err := lists.ListDao.GetList(listId)
	if err != nil {
		return err
	}

	actualConfigs, err := share.ShareConfigDao.GetAllShareConfigsByList(listId)
	if err != nil {
		if err.Status() != http.StatusNotFound {
			return err
		}
	}

	if err := list.ValidateAddItems(callerId, actualConfigs); err != nil {
		return err
	}

	return items.ItemListDao.DeleteItemFromList(itemId, listId)
}

func (l listsService) CheckItem(itemId string, listId int64, callerId int64) apierrors.ApiError {
	list, err := lists.ListDao.GetList(listId)
	if err != nil {
		return err
	}

	actualConfigs, err := share.ShareConfigDao.GetAllShareConfigsByList(listId)
	if err != nil {
		if err.Status() != http.StatusNotFound {
			return err
		}
	}

	if err := list.ValidateCheckItems(callerId, actualConfigs); err != nil {
		return err
	}

	return items.ItemListDao.UpdateItemStatus(itemId, listId, items.StatusChecked)
}
