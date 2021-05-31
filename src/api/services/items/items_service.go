package items_service

import (
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	"github.com/lmurature/melist-api/src/api/domain/items"
	items_provider "github.com/lmurature/melist-api/src/api/providers/items"
)

type itemsService struct{}

type itemsServiceInterface interface {
	SearchItems(query string) (*items.ItemSearchResponse, apierrors.ApiError)
	GetItem(itemId string) (*items.Item, apierrors.ApiError)
}

var (
	ItemsService itemsServiceInterface
)

func init() {
	ItemsService = &itemsService{}
}

func (s *itemsService) SearchItems(query string) (*items.ItemSearchResponse, apierrors.ApiError) {
	result, err := items_provider.SearchItemsByQuery(query)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *itemsService) GetItem(itemId string) (*items.Item, apierrors.ApiError) {
	item, err := items_provider.GetItemById(itemId)
	if err != nil {
		return nil, err
	}

	description, err := items_provider.GetItemDescription(itemId)
	if err != nil {
		return nil, err
	}

	item.Description = description.PlainText

	return item, nil
}
