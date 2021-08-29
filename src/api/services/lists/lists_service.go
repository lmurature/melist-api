package lists

import (
	"fmt"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	"github.com/lmurature/melist-api/src/api/domain/items"
	"github.com/lmurature/melist-api/src/api/domain/lists"
	"github.com/lmurature/melist-api/src/api/domain/notifications"
	"github.com/lmurature/melist-api/src/api/domain/share"
	items_service "github.com/lmurature/melist-api/src/api/services/items"
	users_service "github.com/lmurature/melist-api/src/api/services/users"
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
	GiveAccessToUsers(listId int64, callerId int64, config share.ShareConfigs) (share.ShareConfigs, apierrors.ApiError)
	SearchPublicLists() (lists.Lists, apierrors.ApiError)
	GetMyLists(ownerId int64) (lists.Lists, apierrors.ApiError)
	GetMySharedLists(userId int64) (lists.Lists, apierrors.ApiError)
	AddItemToList(itemId string, variationId int64, listId int64, callerId int64) apierrors.ApiError
	GetItemsFromList(listId int64, callerId int64, info bool) (items.ItemListCollection, apierrors.ApiError)
	DeleteItemFromList(itemId string, listId int64, callerId int64) (items.ItemListCollection, apierrors.ApiError)
	CheckItem(itemId string, listId int64, callerId int64) (items.ItemListCollection, apierrors.ApiError)
	UncheckItem(itemId string, listId int64, callerId int64) (items.ItemListCollection, apierrors.ApiError)
	GetUserFavoriteLists(userId int64) (lists.Lists, apierrors.ApiError)
	MakeFavoriteList(listId int64, userId int64) apierrors.ApiError
	RemoveFavoriteList(listId int64, userId int64) apierrors.ApiError
	GetUserPermissions(listId int64, callerId int64) (*share.ShareConfig, apierrors.ApiError)
	RevokeAccessToUser(listId int64, callerId int64, userId int64) (share.ShareConfigs, apierrors.ApiError)
	GetListNotifications(listId int64, callerId int64) ([]notifications.Notification, apierrors.ApiError)
	GetListItemStatus(itemId string, listId int64, callerId int64) (*items.ItemListDto, apierrors.ApiError)
	GetAllLists() (lists.Lists, apierrors.ApiError)
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

func (l listsService) GiveAccessToUsers(listId int64, callerId int64, config share.ShareConfigs) (share.ShareConfigs, apierrors.ApiError) {
	list, err := lists.ListDao.GetList(listId)
	if err != nil {
		return nil, err
	}

	if err := list.ValidateUpdatability(callerId); err != nil {
		return nil, err
	}

	errorCauseList := make(apierrors.CauseList, 0)
	for i, c := range config {
		if err := c.Validate(); err != nil {
			errorCauseList = append(errorCauseList, err.Message())
		}
		config[i].ListId = listId
	}

	if len(errorCauseList) > 0 {
		return nil, apierrors.NewApiError("invalid request", "invalid request for sharing access to users", http.StatusBadRequest, errorCauseList)
	}

	actualConfigs, err := share.ShareConfigDao.GetAllShareConfigsByList(listId)
	if err != nil {
		if err.Status() != http.StatusNotFound {
			return nil, err
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
		return nil, apierrors.NewApiError("error while trying to save share configs", "database error", http.StatusInternalServerError, errorSaving)
	}

	updatedConfigs, err := share.ShareConfigDao.GetAllShareConfigsByList(listId)
	if err != nil {
		if err.Status() != http.StatusNotFound {
			return nil, err
		}
	}

	return updatedConfigs, nil
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
	items.ItemDao.InsertItem(itemId)

	itemListDto := items.ItemListDto{
		ItemId:      itemId,
		ListId:      listId,
		Status:      items.StatusNotChecked,
		VariationId: variationId,
		UserId:      callerId,
	}

	_, err = items.ItemListDao.InsertItemToList(itemListDto)
	if err != nil {
		return err
	}

	userData, err := users_service.UsersService.GetMeliUser(callerId)
	if err != nil {
		return err
	}

	result, err := notifications.NotificationsDao.SaveNotification(*notifications.NewAddedItemToListNotification(listId, itemId, userData.Nickname))
	if err == nil {
		logrus.Info(fmt.Sprintf("successfully notificated item %s on list %d (%v)", itemId, listId, result))
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
				item, err := items_service.ItemsService.GetItemWithDescription(id)
				output <- items.ItemConcurrent{
					Item:      item,
					ItemError: err,
					ListIndex: index,
				}
			}(itemListCollection[i].ItemId, i, input)
		}

		for i := 0; i < len(itemListCollection); i++ {
			result := <-input
			if result.ItemError != nil {
				return nil, result.ItemError
			}

			itemListCollection[result.ListIndex].MeliItem = result.Item
		}
	}

	return itemListCollection, nil
}

func (l listsService) DeleteItemFromList(itemId string, listId int64, callerId int64) (items.ItemListCollection, apierrors.ApiError) {
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

	if err := list.ValidateAddItems(callerId, actualConfigs); err != nil {
		return nil, err
	}

	err = items.ItemListDao.DeleteItemFromList(itemId, listId)
	if err != nil {
		return nil, err
	}

	return l.GetItemsFromList(listId, callerId, true)
}

func (l listsService) CheckItem(itemId string, listId int64, callerId int64) (items.ItemListCollection, apierrors.ApiError) {
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

	if err := list.ValidateCheckItems(callerId, actualConfigs); err != nil {
		return nil, err
	}

	if err := items.ItemListDao.UpdateItemStatus(itemId, listId, items.StatusChecked); err != nil {
		return nil, err
	}

	userData, err := users_service.UsersService.GetMeliUser(callerId)
	if err != nil {
		return nil, err
	}

	result, err := notifications.NotificationsDao.SaveNotification(*notifications.NewCheckedItemNotification(listId, itemId, userData.Nickname))
	if err == nil {
		logrus.Info(fmt.Sprintf("successfully notificated checked item %s on list %d (%v)", itemId, listId, result))
	}

	return l.GetItemsFromList(listId, callerId, true)
}

func (l listsService) UncheckItem(itemId string, listId int64, callerId int64) (items.ItemListCollection, apierrors.ApiError) {
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

	if err := list.ValidateCheckItems(callerId, actualConfigs); err != nil {
		return nil, err
	}

	if err := items.ItemListDao.UpdateItemStatus(itemId, listId, items.StatusNotChecked); err != nil {
		return nil, err
	}

	userData, err := users_service.UsersService.GetMeliUser(callerId)
	if err != nil {
		return nil, err
	}

	result, err := notifications.NotificationsDao.SaveNotification(*notifications.NewUncheckedItemNotification(listId, itemId, userData.Nickname))
	if err == nil {
		logrus.Info(fmt.Sprintf("successfully notificated unchecked item %s on list %d (%v)", itemId, listId, result))
	}

	return l.GetItemsFromList(listId, callerId, true)
}

func (l listsService) GetUserFavoriteLists(userId int64) (lists.Lists, apierrors.ApiError) {
	return lists.ListDao.GetUserFavoriteLists(userId)
}

func (l listsService) MakeFavoriteList(listId int64, userId int64) apierrors.ApiError {
	_, err := lists.ListDao.GetList(listId)
	if err != nil {
		return err
	}

	if err := lists.ListDao.SaveFavoriteList(listId, userId); err != nil {
		return err
	}

	userData, err := users_service.UsersService.GetMeliUser(userId)
	if err != nil {
		return err
	}

	result, err := notifications.NotificationsDao.SaveNotification(*notifications.NewUserAddedListToFavorites(listId, userData.Nickname))
	if err == nil {
		logrus.Info(fmt.Sprintf("successfully notificated favorited list %d (%v)", listId, result))
	}

	return nil
}

func (l listsService) RemoveFavoriteList(listId int64, userId int64) apierrors.ApiError {
	_, err := lists.ListDao.GetList(listId)
	if err != nil {
		return err
	}

	userFavoriteLists, err := lists.ListDao.GetUserFavoriteLists(userId)
	if err != nil {
		return err
	}

	if !userFavoriteLists.ContainsList(listId) {
		return apierrors.NewBadRequestApiError("list id is not in user's favorites")
	}

	return lists.ListDao.RemoveFavoriteList(listId, userId)
}

func (l listsService) GetUserPermissions(listId int64, callerId int64) (*share.ShareConfig, apierrors.ApiError) {
	list, err := lists.ListDao.GetList(listId)
	if err != nil {
		return nil, err
	}

	if list.OwnerId == callerId {
		return &share.ShareConfig{
			ListId:    listId,
			UserId:    callerId,
			ShareType: share.ShareTypeAdmin,
		}, nil
	}

	shareConfigs, err := share.ShareConfigDao.GetAllShareConfigsByList(listId)
	if err != nil {
		if err.Status() != http.StatusNotFound {
			return nil, err
		}
	}

	for _, conf := range shareConfigs {
		if conf.UserId == callerId {
			return &conf, nil
		}
	}

	if list.Privacy == lists.PrivacyTypePublic {
		return &share.ShareConfig{
			UserId:    callerId,
			ListId:    listId,
			ShareType: share.ShareTypeRead,
		}, nil
	} else {
		return nil, apierrors.NewForbiddenApiError("you have no access to this list")
	}
}

func (l listsService) RevokeAccessToUser(listId int64, callerId int64, userId int64) (share.ShareConfigs, apierrors.ApiError) {
	list, err := lists.ListDao.GetList(listId)
	if err != nil {
		return nil, err
	}

	if err := list.ValidateUpdatability(callerId); err != nil {
		return nil, err
	}

	shareConfigs, err := share.ShareConfigDao.GetAllShareConfigsByList(listId)
	if err != nil {
		return nil, err
	}

	if !slice.ShareConfigUserExists(shareConfigs, userId) {
		return nil, apierrors.NewBadRequestApiError("user does not have access to list")
	}

	var deletedIndex int
	for i, conf := range shareConfigs {
		if conf.UserId == userId {
			deletedIndex = i
			if err := share.ShareConfigDao.DeleteShareConfig(userId, listId); err != nil {
				return nil, err
			}
			break
		}
	}

	return append(shareConfigs[:deletedIndex], shareConfigs[deletedIndex+1:]...), nil
}

func (l listsService) GetListNotifications(listId int64, callerId int64) ([]notifications.Notification, apierrors.ApiError) {
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

	return notifications.NotificationsDao.GetListNotifications(listId)
}

func (l listsService) GetListItemStatus(itemId string, listId int64, callerId int64) (*items.ItemListDto, apierrors.ApiError) {
	listItems, err := l.GetItemsFromList(listId, callerId, false)
	if err != nil {
		return nil, err
	}

	for _, li := range listItems {
		if li.ItemId == itemId {
			return &li, nil
		}
	}

	return nil, apierrors.NewNotFoundApiError(fmt.Sprintf("item %s not found in list %d", itemId, listId))
}

func (l listsService) GetAllLists() (lists.Lists, apierrors.ApiError) {
	return lists.ListDao.GetAllLists()
}
