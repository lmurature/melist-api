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
	var meliItem *items.Item
	var desc *items.ItemDescription
	input := make(chan items.ItemDescriptionConcurrent, 2)
	defer close(input)

	go func(itemId string, output chan items.ItemDescriptionConcurrent) {
		item, err := items_provider.GetItemById(itemId)
		output <- items.ItemDescriptionConcurrent{
			Item:        item,
			Description: nil,
			Error:       err,
		}
	}(itemId, input)

	go func(itemId string, output chan items.ItemDescriptionConcurrent) {
		description, err := items_provider.GetItemDescription(itemId)
		output <- items.ItemDescriptionConcurrent{
			Item:        nil,
			Description: description,
			Error:       err,
		}
	}(itemId, input)

	for i := 0; i < 2; i++ {
		result := <- input
		if result.Error != nil {
			return nil, result.Error
		}

		if result.Item != nil {
			meliItem = result.Item
		} else if result.Description != nil {
			desc = result.Description
		}
	}

	meliItem.Description = desc.PlainText
	return meliItem, nil
}
